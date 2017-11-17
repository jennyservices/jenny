// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package httptestutil

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-openapi/inflect"
)

// ToEnglish given a  *http.Request, *http.Response pair returns a
// name thats understandable in english.
func ToEnglish(req *http.Request, resp *http.Response) string {
	frags := []string{}
	// methods
	{
		action := inflect.Typeify(strings.ToLower(req.Method))
		preposition := "From"
		switch req.Method {
		case http.MethodGet:
			preposition = ""
		case http.MethodPost:
			preposition = "To"
		case http.MethodDelete:
			preposition = "To"
		}
		frags = append(frags, action)

		frags = append(frags, preposition)
	}

	// path
	{
		path := inflect.Camelize(strings.Replace(req.URL.Path, "/", "_", -1))
		frags = append(frags, path)
		ises := []string{}
		for k, v := range req.URL.Query() {
			ises = append(ises, fmt.Sprintf("%sIs%s", inflect.Typeify(k), inflect.Typeify(strings.Join(v, ""))))
		}
		sort.Slice(ises, func(i int, j int) bool {
			return strings.Compare(ises[i], ises[j]) < 0
		})
		if len(ises) > 0 {
			frags = append(frags, "When")
			frags = append(frags, strings.Join(ises, "And"))
		}

	}

	// returning values
	{
		frags = append(frags, "Returns")
		frags = append(frags, inflect.Typeify(http.StatusText(resp.StatusCode)))
	}
	return strings.Join(frags, "")
}
