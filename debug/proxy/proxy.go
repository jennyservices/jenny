// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/jennyservices/jenny/debug/watcher"
	"github.com/cenkalti/backoff"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

var (
	l = log.New(os.Stdout, "proxy: ", log.Lshortfile)
)

type BuildProxy interface {
	http.Handler
	RunLog() []byte
	BuildLog() []byte
	Ready() bool
	Healthy() bool
	Listens() *url.URL
	LivenessPath() string
	SetMiddleware(m Middleware)
}

type Middleware func(http.Handler) http.Handler

type basicProxy struct {
	proxy                  *httputil.ReverseProxy
	updated                time.Time
	cmd                    *exec.Cmd
	buildBuffer, runBuffer *bytes.Buffer
	buildLog               []byte
	middleware             Middleware
	pkgName, bin           string
	pkg                    *build.Package
	u                      *url.URL
	liveness               string
}

func New(pkgName string, middleware Middleware) (BuildProxy, error) {
	// Import the package
	pkg, err := build.Import(pkgName, filepath.Join(build.Default.GOPATH, "src"), build.AllowBinary)
	if err != nil {
		return nil, errors.New("cannot find the go package")
	}
	// Check if the package is an executable
	if !pkg.IsCommand() {
		return nil, errors.New("is command: serve requires a command that allows -addr")
	}

	_, bin := path.Split(pkg.ImportPath)

	bp := &basicProxy{
		updated:     time.Time{},
		proxy:       nil,
		cmd:         nil,
		pkgName:     pkgName,
		pkg:         pkg,
		bin:         bin,
		buildBuffer: &bytes.Buffer{},
		runBuffer:   &bytes.Buffer{},
		middleware:  middleware,
		u:           nil,
		liveness:    "/_ready",
	}
	go bp.start(context.Background())
	return bp, nil
}

func (p *basicProxy) LivenessPath() string {
	return p.liveness
}

func (p *basicProxy) SetMiddleware(m Middleware) {
	p.middleware = m
}
func getID(r *http.Request) string {
	reqID := r.Header.Get("X-Request-Id")
	if reqID != "" {
		return reqID
	}
	id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}

func (b *basicProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if b.proxy == nil {
		io.WriteString(w, "application not ready")
		return
	}
	if b.middleware != nil {
		b.middleware(b.proxy).ServeHTTP(w, r)
	} else {
		b.proxy.ServeHTTP(w, r)
	}
}

func (b *basicProxy) Ready() bool {
	return b.cmd != nil && b.u != nil
}

func (b *basicProxy) Healthy() bool {
	return true
}

func (b *basicProxy) Listens() *url.URL {
	return b.u
}

func (b *basicProxy) BuildLog() []byte {
	return b.buildBuffer.Bytes()
}

func (b *basicProxy) RunLog() []byte {
	return b.runBuffer.Bytes()
}

func (b *basicProxy) needsRefresh() bool {
	w := watcher.NewPollingPackageWatcher(b.pkgName, b.updated)
	if w.Updated() {
		b.updated = time.Now()
		return true
	}
	return false
}

// Get - return free open TCP port
func get() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return 0, err
	}
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	time.Sleep(time.Millisecond * 100)
	return
}

func (b *basicProxy) runNewProxy(ctx context.Context, cmd string) error {

	rl := log.New(b.runBuffer, "runner: ", log.Lshortfile)
	port, err := get()
	if err != nil {
		rl.Fatal(err)
	}
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	args := []string{
		"-addr",
		addr,
	}
	rl.Println(cmd, args)
	runCmd := exec.CommandContext(ctx, cmd, args...)
	runCmd.Stderr = b.runBuffer
	runCmd.Stdout = b.runBuffer
	if err = runCmd.Start(); err != nil {
		rl.Println(errors.Wrap(err, "run start"))

		return err
	}

	time.Sleep(time.Millisecond * 100)
	if runCmd.Process != nil {
		rl.Println("started ", cmd)
	} else {
		return fmt.Errorf("process has died")
	}
	baseURL := "http://" + addr
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	u.Path = b.liveness
	operation := func() error {
		resp, err := http.DefaultClient.Get(u.String())
		if err != nil {
			l.Println(err)
			return backoff.Permanent(err)
		}
		switch resp.StatusCode {
		case http.StatusServiceUnavailable:
			return errors.New("service not ready yet")
		case http.StatusOK:
			return nil
		default:
			return backoff.Permanent(err)
		}
	}

	if err = backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
		//runCmd.Process.Kill()
		rl.Printf("%q with pid(%d) has started but is not healthy: %v\n", cmd, runCmd.Process.Pid, err)
		return err
	}

	rl.Printf("%s (%d) has started and is healthy\n", cmd, runCmd.Process.Pid)

	go func() {
		err = runCmd.Wait()
		if err != nil {
			rl.Println(err)
		}
		rl.Printf("%s (%d) has quit", cmd, runCmd.Process.Pid)
	}()

	b.proxy = httputil.NewSingleHostReverseProxy(u)
	b.u = u
	if b.cmd != nil {
		x := b.cmd
		if err = x.Process.Signal(os.Interrupt); err != nil {
			return err
		}
		err = b.cmd.Wait()
		if err != nil {
			return err
		}
	}

	b.cmd = runCmd

	return nil
}

func (b *basicProxy) build(ctx context.Context) (string, error) {
	bl := log.New(b.buildBuffer, "build: ", log.Lshortfile)
	tmpdir, err := ioutil.TempDir("", "jenny")
	if err != nil {
		bl.Fatal(err)
	}

	binary := path.Join(tmpdir, b.bin)
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binary)
	cmd.Dir = b.pkg.Dir
	cmd.Stderr = b.buildBuffer
	cmd.Stdout = b.buildBuffer

	bl.Println("go", "build", "-o", binary)
	err = cmd.Run()
	if err != nil {
		bl.Println("compilation failed")
	} else {
		bl.Println("compiled", b.bin)
	}
	return binary, err
}

func (b *basicProxy) start(ctx context.Context) {

	go http.ListenAndServe(":8080", b)
	for {
		if b.needsRefresh() {
			bin, err := b.build(ctx)
			if err == nil {
				l.Println("compiled", bin)
				if err := b.runNewProxy(ctx, bin); err != nil {
					l.Println(errors.Wrap(err, "proxy reloader"))
				}
			} else {
				l.Println(err)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}
