// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jenny

import (
	"io"

	"github.com/Typeform/jenny/generator/internal/ir"

	"github.com/Typeform/jenny/generator/swagger"
)

type Service ir.Service
type Method ir.Method

func DecodeService(r io.Reader) (*Service, error) {
	dec := swagger.NewDecoder(r)
	svc, err := dec.Decode()
	x := Service(*svc)
	return &x, err
}
