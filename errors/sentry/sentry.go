// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package sentry implements the Reporter interface defined in errors for sentry
// see sentry at https://sentry.io/welcome/
package sentry

import (
	"context"
	"os"
	"runtime"

	pkgErrors "github.com/pkg/errors"

	"github.com/Typeform/jenny/errors"
	"github.com/Typeform/jenny/http"
	raven "github.com/getsentry/raven-go"
)

type sentry struct {
	hostname, release string
	rvn               *raven.Client
}

// New returns a new reporter that reports errors to Sentry
func New(dsn, release string) (errors.Reporter, error) {
	c, err := raven.New(dsn)
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &sentry{
		rvn:      c,
		hostname: hostname,
		release:  release,
	}, nil
}

// Report reports an error to sentry
func (s *sentry) Report(ctx context.Context, err error, op string) {
	id := string(http.ContextRequestID(ctx))
	cause := pkgErrors.Cause(err)
	pkt := raven.NewPacket(cause.Error(), raven.NewException(cause, raven.GetOrNewStacktrace(cause, 1, 3, []string{})))
	extra := map[string]interface{}{
		"operation":    op,
		"x_request_id": id,
		"runtime": map[string]interface{}{
			"version":       runtime.Version(),
			"num_cpu":       runtime.NumCPU(),
			"max_procs":     runtime.GOMAXPROCS(0),
			"num_goroutine": runtime.NumGoroutine(),
		},
	}
	pkt.Extra = extra
	s.rvn.Capture(pkt, nil)
}
