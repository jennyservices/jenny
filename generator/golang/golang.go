// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package golang

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Typeform/jenny/generator/internal/ir"
	"github.com/Typeform/jenny/generator/util"
	"github.com/go-openapi/inflect"
	"github.com/pkg/errors"
	"sevki.org/lib/source"
)

const (
	str = "string"
	obj = "object"
)

var (
	base  = path.Join(build.Default.GOPATH, "src/github.com/Typeform/jenny/templates/go/")
	tmpl  = template.Must(template.New("").Funcs(funcs).ParseGlob(path.Join(base, "/*.go.tmpl")))
	pathy = path.Join(build.Default.GOPATH, "src/github.com/Typeform/jenny/templates/go/*.go.tmpl")
	fset  = token.NewFileSet()

	funcs = map[string]interface{}{
		"join":              strings.Join,
		"camelizeDownFirst": camelizeDownFirst,
		"camelize":          inflect.Camelize,
		"normalize":         util.NormalizeName,
		"titleize":          inflect.Titleize,
		"getType":           getType,
		"firstOrDefault":    firstOrDefault,
		"isFromLocation": func(p *string, s string) bool {
			if p != nil && s == *p {
				return true
			}
			return false
		},
	}
)

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

func parseType(s string) string {
	isArray := false
	if strings.HasPrefix(s, "[]") {
		isArray = true
		s = strings.TrimLeft(s, "[]")
	}
	_, f := path.Split(s)
	if isArray {
		return "[]" + f
	}
	return f
}

func getType(schema ir.Schema) string {
	s := getSimpleType(schema.Type, schema.Format)
	isArray := false
	if strings.HasPrefix(s, "[]") {
		isArray = true
	}
	if schema.Required || isArray {
		return s
	}
	return "*" + s
}

func getSimpleType(simpleType string, f *string) string {
	format := ""
	if f != nil {
		format = *f
	}
	switch simpleType {
	case "number":
		switch format {
		case "float", "double":
			return "float64"
		default:
			return "int"
		}
	case "integer":
		switch format {
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		default:
			return "int"
		}
	case str:
		switch format {
		case "date", "date-time":
			return "time.Time"
		case "byte", "binary":
			return "[]byte"
		case "password":
			return str
		case "url":
			return "url.URL"
		default:
			return str
		}
	case "boolean":
		return "bool"
	default:
		return parseType(simpleType)
	}
}

// New returns a new goWriter
func New(w io.Writer, template string) ir.Encoder {
	return &goWriter{
		template: template,
		w:        w,
	}
}

func String(s string) *string { x := s; return &x }

// Encode encodes a service defintion to a go file
func (g *goWriter) Encode(s *ir.Service) error {
	buf := bytes.NewBuffer(nil)

	if err := tmpl.ExecuteTemplate(buf, g.template, s); err != nil {
		return err
	}

	bytz := buf.Bytes()
	fmtd, err := processFile(buf)
	if err != nil {
		formattingError(err, bytz)
		return errors.Wrap(err, "formatting")
	}

	g.w.Write(fmtd)
	return nil
}

func formattingError(err error, bytz []byte) {
	lines := source.GetRangesFromErrors(source.ParseSourceErrors(err.Error()), 10)
	scanner := bufio.NewScanner(bytes.NewBuffer(bytz))
	x := make(map[int]string)
	i := 1
	for scanner.Scan() {
		x[i] = scanner.Text()
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	for _, line := range lines {
		if text, ok := x[line]; ok {
			log.Printf("%d: %s\n", line, text)
		}
	}
}

func camelizeDownFirst(s string) string {
	if s == "" {
		return ""
	}
	return inflect.CamelizeDownFirst(s)
}

type goWriter struct {
	w        io.Writer
	template string
}
