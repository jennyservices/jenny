// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package watcher

import (
	"go/build"
	"os"
	"path/filepath"
	"time"
)

const src = "src"

// Watcher represents a inferface that lets you know if there are changes
// in a filesystem you are looking at.
type Watcher interface {
	Updated() bool
}

type pollingPackageWatcher struct {
	latest  time.Time
	seen    map[string]bool
	pkgName string
}

// NewPollingPackageWatcher tracks if there are changes in a go package
// that would warrant recompilation. It is not as sophisticated as bldy so
// it is suseptible to racy
func NewPollingPackageWatcher(pkgName string, t time.Time) Watcher {
	return &pollingPackageWatcher{
		latest:  t,
		seen:    make(map[string]bool),
		pkgName: pkgName,
	}
}

func (p *pollingPackageWatcher) Updated() bool {
	x := p.allPkgFiles(p.pkgName)
	for _, file := range x {
		info, err := os.Stat(file)
		if err == nil && info.ModTime().After(p.latest) {
			return true
		}
	}
	return false
}

func (p *pollingPackageWatcher) allPkgFiles(pkgName string) (x []string) {
	if _, seen := p.seen[pkgName]; seen {
		return nil
	}
	p.seen[pkgName] = true

	// Import the package
	pkg, err := build.Import(pkgName, filepath.Join(build.Default.GOPATH, src), build.AllowBinary)
	if err != nil {
		return nil
	}
	for _, file := range pkg.GoFiles {
		x = append(x, filepath.Join(pkg.Dir, file))
	}
	for _, imprt := range pkg.Imports {
		x = append(x, p.allPkgFiles(imprt)...)
	}
	return x
}
