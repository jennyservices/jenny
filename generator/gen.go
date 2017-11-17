// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Typeform/jenny/generator/golang"
	"github.com/Typeform/jenny/generator/internal/ir"
	"github.com/Typeform/jenny/generator/js"
	"github.com/Typeform/jenny/generator/swagger"

	"github.com/pkg/errors"
)

type encoder int
type decoder int

const (
	Swagger decoder = iota
	Golang  encoder = iota
	JS
)

// Generator is an object that can be used to create jenny services
type Generator struct {
	out, in    string
	Package    string
	filePrefix string
}

// New returns a new generator with the given output directory
func New(in, out string) *Generator {
	return &Generator{
		in:         in,
		out:        out,
		filePrefix: "jenny",
	}
}

func (g *Generator) SetPrefix(s string) {
	g.filePrefix = s
}

// Generate Generates the necessary folders and files for the typeform service.
func (g *Generator) Generate() error {
	dec, err := swagger.NewDecoderFromFile(g.in)
	if err != nil {
		log.Fatal(err)
	}

	svc, err := dec.Decode()
	if err != nil {
		return errors.Wrap(err, "generate")
	}

	for _, output := range map[string]struct {
		file     string
		encoder  encoder
		template string
	}{
		"service": {
			file:     fmt.Sprintf("%s_service.go", g.filePrefix),
			encoder:  Golang,
			template: "servicebare",
		},
		"http": {
			file:     fmt.Sprintf("%s_http.go", g.filePrefix),
			encoder:  Golang,
			template: "servicehttp",
		},
	} {
		f, err := os.Create(filepath.Join(g.out, output.file))
		if err != nil {
			return errors.Wrap(err, "create file")
		}

		var enc ir.Encoder
		switch output.encoder {
		case Golang:
			enc = golang.New(f, output.template)
		case JS:
			enc = js.New(f, output.template)
		}

		if err := enc.Encode(svc); err != nil {
			return errors.Wrap(err, "encoder")
		}
	}

	return nil
}
