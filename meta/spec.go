// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package spec handles the service definitions
package spec

import (
	"context"
	"crypto/sha256"
	"net/http"

	"github.com/jennyservices/jenny/spec/transport/v1"
	swagger "github.com/go-openapi/spec"
)

var (
	s       *specsService
	service http.Handler
)

// Register is used to register a swagger spec under a handle
func Register(name string, spec []byte) {
	if s == nil {
		s = &specsService{
			specs: make(map[string][]byte),
		}
		service = v1.NewSpecsService(s)
	}
	s.specs[name] = spec
}

type specsService struct {
	specs map[string][]byte
}

func (s *specsService) Specs(ctx context.Context) ([]*v1.Spec, error) {
	x := []*v1.Spec{}
	for name, spec := range s.specs {
		doc := &swagger.Swagger{}
		h := sha256.New()
		h.Write(spec)
		doc.UnmarshalJSON(spec)
		x = append(x, &v1.Spec{
			Name:              name,
			SwaggerDefinition: spec,
			Version:           doc.Info.Version,
			Hash:              h.Sum(nil),
		})
	}
	return x, nil
}

// Handler returns the handler for the service
func Handler() http.Handler {
	return service
}
