// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jennyservices/jenny/debug/proxy"
	"github.com/jennyservices/jenny/debug/transport/v1"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type debugService struct {
	pkg    string
	db     *bolt.DB
	wd     string
	runBuf *bytes.Buffer
	trips  map[string]trip
	proc   *os.Process
	proxy  proxy.BuildProxy
}

// Start starts the debug server.
// Start Requires go pkgs to start properly.
func Start(pkgName string) {
	// Create a new jennydb
	db, err := bolt.Open(".jennydebug.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	p, err := proxy.New(pkgName, nil)
	if err != nil {
		log.Fatal(err)
	}

	debugServer := &debugService{
		pkg:   pkgName,
		db:    db,
		proxy: p,
		trips: make(map[string]trip),
	}

	p.SetMiddleware(debugServer.debugHandler)

	ds := v1.NewDebugHTTPServer(debugServer).(*mux.Router)
	jennyUI := filepath.Join(os.Getenv("GOPATH"), "src/github.com/jennyservices/jenny/debug/transport/ui/")
	ds.PathPrefix("/").Handler(http.FileServer(http.Dir(jennyUI)))

	log.Println("starting debug server at http://0.0.0.0:1337")
	log.Fatal(http.ListenAndServe(":1337", cors.Default().Handler(ds)))
}
