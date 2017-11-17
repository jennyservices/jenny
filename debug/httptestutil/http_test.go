// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package httptestutil

import (
	"net/http"
	"net/url"
	"testing"
)

var httpUtilTests = []struct {
	name     string
	req      *http.Request
	resp     *http.Response
	expected string
}{
	{
		req: &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Path: "/home/",
			},
		},
		resp: &http.Response{
			StatusCode: 200,
		},
		expected: "GetHomeReturnsOK",
	},
	{
		req: &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Path: "/home",
			},
		},
		resp: &http.Response{
			StatusCode: 200,
		},
		expected: "GetHomeReturnsOK",
	},
	{
		req: &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Path:     "/home",
				RawQuery: "id=1",
			},
		},
		resp: &http.Response{
			StatusCode: 200,
		},
		expected: "GetHomeWhenIdIs1ReturnsOK",
	},
	{
		req: &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Path:     "/home",
				RawQuery: "id=1&sort=desc",
			},
		},
		resp: &http.Response{
			StatusCode: 200,
		},
		expected: "GetHomeWhenIdIs1AndSortIsDescReturnsOK",
	},
	{
		req: &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Path:     "/home",
				RawQuery: "id=1&sort=desc",
			},
		},
		resp: &http.Response{
			StatusCode: 200,
		},
		expected: "PostToHomeWhenIdIs1AndSortIsDescReturnsOK",
	},
	{
		req: &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Path:     "/home",
				RawQuery: "id=1&sort=desc",
			},
		},
		resp: &http.Response{
			StatusCode: 201,
		},
		expected: "PostToHomeWhenIdIs1AndSortIsDescReturnsCreated",
	},
	{
		req: &http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				Path: "/home",
			},
		},
		resp: &http.Response{
			StatusCode: 200,
		},
		expected: "DeleteToHomeReturnsOK",
	},
}

func TestToEnglish(t *testing.T) {
	for _, test := range httpUtilTests {
		t.Run(test.name, func(t *testing.T) {
			actual := ToEnglish(test.req, test.resp)
			if test.expected != actual {
				t.Logf("was expecting %q got %q instead", test.expected, actual)
				t.Fail()
			}
		})
	}
}
