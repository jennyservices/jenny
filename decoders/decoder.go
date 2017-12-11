// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package decoders

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Typeform/jenny/mime"
	"github.com/gorilla/schema"
)

var (
	// ErrDecoderNotFound is returned when a reuqest doesn't have
	// enough information to determine a decoder
	ErrDecoderNotFound = errors.New("decoder could not be found")
	// JSONDecoder decodes data from a http.Request
	JSONDecoder = func(r *http.Request) Decoder {
		return json.NewDecoder(r.Body)
	}
	// XMLDecoder decodes data from a http.Request
	XMLDecoder = func(r *http.Request) Decoder {
		return xml.NewDecoder(r.Body)
	}
	// FormDecoder decodes data from a http.Request
	FormDecoder = func(r *http.Request) Decoder {
		return &formDecoder{r: r}
	}
	decoders = map[mime.Type]newDecoder{
		mime.ApplicationJSON:           JSONDecoder,
		mime.ApplicationXML:            XMLDecoder,
		mime.ApplicationFormURLEncoded: FormDecoder,
	}
)

// Register registers a new decoder to be used with jenny endpoints, it is to be
// recalled based on the mime-type
func Register(s mime.Type, d newDecoder) {
	decoders[s] = d
}

// Decoder is an interface that decodes http.Request.Body from their
// Content-Type mime types.
type Decoder interface {
	Decode(v interface{}) error
}

type newDecoder func(*http.Request) Decoder

type formDecoder struct {
	r *http.Request
}

// TokenRequest as defined in swagger
type TokenRequest struct {
	ClientID  string `json:"client_id" schema:"client_id"`
	Code      string `json:"code" schema:"code"`
	GrantType string `json:"grant_type" schema:"grant_type"`
}

func (f *formDecoder) Decode(i interface{}) error {
	dec := schema.NewDecoder()
	body, err := ioutil.ReadAll(f.r.Body)
	log.Println(string(body))
	if err != nil {
		return fmt.Errorf("decoding form: reading body: %v", err)
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return fmt.Errorf("decoding form: parsing values: %v", err)
	}
	if len(values) <= 0 {
		return fmt.Errorf("decoding form: no values found")
	}
	return dec.Decode(i, values)
}

// RequestDecoder returns a decoder for a given http.Request
func RequestDecoder(r *http.Request, accepts []mime.Type) (Decoder, error) {
	serverAccepts := mime.Aggregate(accepts)
	clientSent := mime.NewTypes(r.Header.Get("Content-Type"))
	available := mime.Intersect(serverAccepts, clientSent)
	var dec Decoder
	err := available.Walk(func(x mime.Type) error {
		if decoderFunc, ok := decoders[x]; ok {
			dec = decoderFunc(r)
			return nil
		}
		return fmt.Errorf("%s isn't a registered decoder", x)

	})
	if err != nil {
		return nil, err
	}
	if dec == nil {
		return nil, errors.New("coudln't find decoder")
	}
	return dec, nil
}
