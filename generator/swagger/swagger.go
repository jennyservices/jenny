// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package swagger

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/jennyservices/jenny/generator/internal/ir"
	"github.com/go-openapi/inflect"

	"github.com/blang/semver"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"sevki.org/x/errors"
)

var (
	// ErrSpecDoesntExist is returned when a file doesn't exist
	ErrSpecDoesntExist = errors.New("swagger: spec doesn't exist")
	// ErrSpecIncorrect is returned when a spec is incorrect
	ErrSpecIncorrect = errors.New("swagger: spec is not correct")

	isTest = false
)

type swaggerDecoder struct {
	r    io.Reader
	svc  *ir.Service
	spec *spec.Swagger
}

// NewDecoder returns a new generator with the given output directory
func NewDecoder(r io.Reader) ir.Decoder {
	return &swaggerDecoder{
		r: r,
	}
}

// NewDecoderFromFile takes a swagger 2.0 file and creates a swagger.Decoder
func NewDecoderFromFile(file string) (ir.Decoder, error) {
	stat, err := os.Stat(file)
	if err != nil || stat == nil {
		return nil, ErrSpecDoesntExist
	}

	doc, err := loads.Spec(file)
	if err != nil {
		return nil, ErrSpecIncorrect
	}
	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// doc, err = doc.Expanded(&spec.ExpandOptions{
	// 	ContinueOnError: false,
	// 	SkipSchemas:     false,
	// 	RelativeBase:    wd,
	// })

	if err != nil {
		return nil, err
	}

	specBytz, _ := doc.Spec().MarshalJSON()
	if doc.Spec().Info == nil {
		return nil, ErrSpecIncorrect
	}
	specBuf := bytes.NewBuffer(specBytz)
	return NewDecoder(specBuf), nil
}

// Deecode decodes the contents of a swagger definition to an Intermediary
// Service Definition that is used by Jenny. This Intermediary service
// definition is not to be written by hand but rather converted in to by other
// service definition formats like swagger {2.0 and 3.0}, protobuf and so on.
func (s *swaggerDecoder) Decode() (*ir.Service, error) {

	bytz, err := ioutil.ReadAll(s.r)
	if err != nil {
		return nil, ErrSpecDoesntExist
	}

	s.spec = &spec.Swagger{}
	s.spec.UnmarshalJSON(bytz)

	if s.spec.Info == nil {
		return nil, ErrSpecIncorrect
	}

	s.svc = &ir.Service{
		Schemas:   make(map[string]ir.Schema),
		Resources: make(map[string]ir.Method),
	}

	s.etag(bytz)

	g := errors.NewGroup("decode")
	g = g.Add(s.decodeMeta())
	g = g.Add(s.decodeSchemas())
	g = g.Add(s.decodeMethods())
	if g.Errored() {
		return nil, g
	}
	return s.svc, nil
}

func (s *swaggerDecoder) etag(bytz []byte) {
	h := sha256.New()
	h.Write(bytz)
	s.svc.Etag = fmt.Sprintf("%x", h.Sum(nil))
}

func (s *swaggerDecoder) decodeMeta() *errors.Group {
	spec := s.spec
	g := errors.NewGroup("meta")

	version, err := semver.New(spec.Info.Version)
	g = g.Add(err)

	name := inflect.Camelize(spec.Info.Title)

	s.svc.BasePath = spec.BasePath
	s.svc.BaseURL = spec.Host
	s.svc.Kind = "jenny.io/service"
	s.svc.ID = fmt.Sprintf("%s:v%d.%d.%d", name, version.Major, version.Minor, version.Patch)
	s.svc.DiscoveryVersion = fmt.Sprintf("v%d", version.Major)
	s.svc.CanonicalName = spec.Info.Title
	s.svc.Name = name
	s.svc.Version = spec.Info.Version
	s.svc.Description = spec.Info.Description
	s.svc.Title = spec.Info.Title
	s.svc.DocumentationLink = "https://localhost:8080/_spec"
	s.svc.Protocol = "http"
	url, err := url.Parse(spec.Host)
	g = g.Add(err)

	if spec.Info.Contact != nil {
		s.svc.OwnerName = spec.Info.Contact.Name
	}

	s.svc.RootURL = url.Host

	revision := "deadbeef"

	if out, err := exec.Command("git", []string{"describe", "--always"}...).Output(); err == nil {
		revision = strings.TrimSpace(string(out))
	}

	if isTest {
		revision = "TEST"
	}

	s.svc.Revision = revision

	return g
}
