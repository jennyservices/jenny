package swagger

import (
	"encoding/json"
	"flag"
	"go/build"
	"io"
	"os"
	"path"
	"testing"

	"github.com/jennyservices/jenny/generator/internal/ir"
	"github.com/d4l3k/messagediff"

	"sevki.org/lib/prettyprint"
)

func init() {
	isTest = true
}

var (
	update = flag.Bool("u", false, "-u")
)

func TestSwagger(t *testing.T) {
	flag.Parse()
	tests := []struct {
		name string
		file string
		err  error
	}{
		{
			name: "petstore-oauth",
			file: "./testdata/petstore-oauth.yaml",
			err:  nil,
		},
		{
			name: "petstore",
			file: "./testdata/petstore.yaml",
			err:  nil,
		},
		{
			name: "debug",
			file: path.Join(build.Default.GOPATH, "src/github.com/jennyservices/jenny/debug/transport/v1/swagger.yaml"),
			err:  nil,
		},
		{
			name: "NOFILE",
			file: "/NON/EXISTENT/FILE",
			err:  ErrSpecDoesntExist,
		},
		{
			name: "WRONGSPEC",
			file: path.Join(build.Default.GOPATH, "src/github.com/jennyservices/jenny/.travis.yml"),
			err:  ErrSpecIncorrect,
		},
		{
			name: "NOTASPEC",
			file: path.Join(build.Default.GOPATH, "src/github.com/jennyservices/jenny/README.md"),
			err:  ErrSpecIncorrect,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			dec, err := NewDecoderFromFile(test.file)
			if err != test.err {
				t.Error(err)
				t.Fail()
				return
			} else if err != nil {
				return
			}

			svc, err := dec.Decode()
			if err != nil && err != test.err {
				t.Log(err)
				t.Fail()
				return
			} else if err != nil {
				return
			}

			gold := "./testdata/" + test.name + ".gold.json"

			var goldFile *os.File
			if *update {
				goldFile, err = os.Create(gold)
				if err != nil {
					t.Fatal(err)
				} // Copyright 2017 Typeform SL. All rights reserved.
				// Use of this source code is governed by a MIT-style
				// license that can be found in the LICENSE file.

				io.WriteString(goldFile, prettyprint.AsJSON(svc))

				goldFile.Close()
				return
			}

			goldFile, err = os.Open(gold)
			if err != nil {
				t.Fail()
				t.Log(err)
			}

			jdec := json.NewDecoder(goldFile)
			goldSvc := &ir.Service{}
			if err := jdec.Decode(goldSvc); err != nil {
				t.Fail()
				t.Log(err)
			}
			if x, ok := messagediff.PrettyDiff(svc, goldSvc); !ok {
				t.Fail()
				t.Log("\n" + x)
			}
		})
	}
}
