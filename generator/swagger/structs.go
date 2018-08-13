// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package swagger

import (
	"bytes"
	"fmt"

	"github.com/jennyservices/jenny/generator/internal/ir"
	"github.com/jennyservices/jenny/generator/util"
	"sevki.org/lib/errors"
	"sevki.org/lib/prettyprint"

	"github.com/go-openapi/spec"
)

const (
	object = "object"
	array  = "array"
)

func (s *swaggerDecoder) decodeSchemas() *errors.Group {
	g := errors.NewGroup("schemas")
	for definitionName, definition := range s.spec.Definitions {
		schema, err := s.swaggerToIR(definitionName, s.svc.Name, &definition)
		g = g.Add(err)

		if !g.Errored() {
			s.svc.Schemas[schema.Name] = *schema
		}
	}
	return g
}

func contains(a []string, s string) bool {
	for _, b := range a {
		if s == b {
			return true
		}
	}
	return false
}

func newID(p, c string) string {
	return fmt.Sprintf("%s/%s", p, util.NormalizeName(c))
}

func (s *swaggerDecoder) describe(x *ir.Schema, schema *spec.Schema) *errors.Group {
	g := errors.NewGroup("describe")
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s is generated from a swagger definition", x.Name))
	if !g.Errored() {
		s := buf.String()
		x.Description = &s
	}
	return g
}

func (s *swaggerDecoder) swaggerToIR(name, parent string, schema *spec.Schema) (*ir.Schema, *errors.Group) {
	x := &ir.Schema{
		Name:           util.NormalizeName(name),
		CannonicalName: name,
		ID:             newID(parent, util.NormalizeName(name)),
	}

	g := errors.NewGroup("swagger to ir")
	g = s.guessType(x, schema)
	g = s.describe(x, schema)

	for name, sc := range schema.Properties {
		prop, err := s.swaggerToIR(name, x.ID, &sc)
		if err != nil && err.Errored() {
			g = g.Add(err)
			continue
		}
		if x.Properties == nil {
			x.Properties = make(map[string]ir.Schema)
		}
		x.Properties[prop.Name] = *prop
	}
	for _, required := range schema.Required {
		if prop, ok := x.Properties[util.NormalizeName(required)]; ok {
			prop.Required = true
			x.Properties[util.NormalizeName(required)] = prop
		}
	}
	if g.Errored() {
		return nil, g
	}
	return x, g
}

func (s *swaggerDecoder) guessType(x *ir.Schema, schema *spec.Schema) *errors.Group {
	g := errors.NewGroup("type")

	typeSchema := s.typer(schema)
	if typeSchema == nil {
		return g.Newf("couldn't determine type for %q", x.ID)
	}
	x.Type = typeSchema.Type
	x.Format = typeSchema.Format
	switch x.Type {
	case object:
		x.Type = x.ID
	case array:
		if schema.Items == nil {
			return g.Newf("items for %q shouldn't be empty", x.ID)
		}
		arrayType := s.typer(schema.Items.Schema)
		if arrayType == nil {
			return g.Newf("coudln't determine array type for %q\n%s", x.ID, prettyprint.AsJSON(schema))
		}
		x.Type = "[]" + arrayType.Type
	}

	return g
}

func (s *swaggerDecoder) typer(schema *spec.Schema) *ir.Type {
	x := &ir.Type{}
	if len([]string(schema.Type)) > 0 {
		x.Type = schema.Type[0]
		if len(schema.Format) > 0 {
			x.Format = String(schema.Format)
		} else {
			x.Format = nil
		}
		return x
	}
	if len(schema.Properties) > 0 {
		x.Type = object
		return x
	}
	if schema.Items != nil {
		x.Type = array
		return x
	}
	if ptr := schema.Ref.GetPointer(); ptr != nil {
		frags := ptr.DecodedTokens()
		if len(frags) < 2 {
			return nil
		}

		p := newID(s.svc.Name, frags[1])
		x.Type = p
		return x
	}

	return nil
}

func String(s string) *string { x := s; return &x }
