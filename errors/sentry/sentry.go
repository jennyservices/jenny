// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package sentry implements the Reporter interface defined in errors for sentry
// see sentry at https://sentry.io/welcome/
package sentry

import (
	"os"

	"github.com/Typeform/jenny/errors"
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
func (s *sentry) Report(err error, op string) {

	s.rvn.CaptureError(err, map[string]string{
		"operation": op,
		"release":   s.release,
		"hostname":  s.hostname,
	})
}
