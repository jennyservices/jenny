// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package js

import (
	"go/build"
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/Typeform/jenny/generator/internal/ir"
	"github.com/go-openapi/inflect"
)

var (
	tmpl  = template.Must(template.New("").Funcs(funcs).ParseGlob(path.Join(build.Default.GOPATH, "src/github.com/Typeform/jenny/templates/js/*.js.tmpl")))
	funcs = map[string]interface{}{
		"join":              strings.Join,
		"camelizeDownFirst": camelizeDownFirst,
		"camelize":          inflect.Camelize,
		"titleize":          inflect.Titleize,
	}
)

// New returns a new goWriter
func New(w io.Writer, template string) ir.Encoder {
	return &jsWriter{
		template: template,
		w:        w,
	}
}

// Encode encodes a service defintion to a go file
func (g *jsWriter) Encode(s *ir.Service) error {
	return tmpl.ExecuteTemplate(g.w, g.template, s)
}

func camelizeDownFirst(s string) string {
	if s == "" {
		return ""
	}
	return inflect.CamelizeDownFirst(s)
}

type jsWriter struct {
	w        io.Writer
	template string
}
