// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package js

import (
	"go/build"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/jennyservices/jenny/generator/internal/ir"
	"github.com/jennyservices/jenny/generator/util"
	"github.com/go-openapi/inflect"
)

var (
	tmpl  = template.Must(template.New("").Funcs(funcs).ParseGlob(path.Join(build.Default.GOPATH, "src/github.com/jennyservices/jenny/templates/js/*.js.tmpl")))
	funcs = map[string]interface{}{
		"join":              strings.Join,
		"camelizeDownFirst": camelizeDownFirst,
		"camelize":          inflect.Camelize,
		"normalize":         util.NormalizeName,
		"titleize":          inflect.Titleize,
		"getType":           getType,
		"firstOrDefault":    firstOrDefault,
		"tableize":          tableUpper,
		"uppertable":        strings.ToUpper,
		"isFromLocation": func(p *string, s string) bool {
			if p != nil && s == *p {
				return true
			}
			return false
		},
	}
)

func tableUpper(s string) string {
	return strings.ToUpper(inflect.Tableize(s))
}
func parseType(s string) string {
	isArray := false
	if strings.HasPrefix(s, "[]") {
		isArray = true
		s = strings.TrimLeft(s, "[]")
	}
	_, f := path.Split(s)
	if isArray {
		return f + "[]"
	}
	return f
}

func getType(schema ir.Schema) string {
	return parseType(schema.Type)
}

func firstOrDefault(m map[string]ir.Response) ir.Response {

	r := ir.Response{
		Default:    true,
		Error:      false,
		HTTPStatus: 200,
		Returns:    map[string]ir.Schema{},
	}
	for _, s := range m {
		r = s
		break
	}
	for _, s := range m {
		if s.Default {
			return s
		}
	}
	return r
}

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
