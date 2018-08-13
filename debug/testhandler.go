// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/jennyservices/jenny/debug/httptestutil"
	"github.com/jennyservices/jenny/debug/transport/v1"
	"golang.org/x/tools/imports"
)

func createTest(trip *v1.Trip) {
	buf := bytes.NewBuffer(nil)

	testTemplate := `
func Test{{.Name}}(t *testing.T){
{{ with .Trip }}

	// Request
	requestURL:= "{{ .Request.URL }}";
	requestMethod := "{{ .Request.Method }}";
	requestHeader := make(http.Header)
	for k, v := range map[string][]string {
		{{- range $k, $v := .Request.Headers  }}
		"{{$v.Key}}": []string{ {{range $v.Values }} "{{ .}}",{{end}} },
		{{- end }}
	}{
		for _, val := range v {
			requestHeader.Add(k,val)
		}
	}

	requestBody := []byte{ {{range .Request.Body }}{{ . }},{{end}} }

	// Response
	responseStatus := {{.Response.Status}}
	responseHeader := make(http.Header)
	for k, v := range map[string][]string {
		{{- range $k, $v := .Response.Headers  }}
		"{{$v.Key}}": []string{ {{range $v.Values }} "{{ .}}",{{end}} },
		{{- end }}
	}{
		for _, val := range v {
			requestHeader.Add(k,val)
		}
	}
	responseBody := []byte{ {{range .Response.Body }}{{ . }},{{end}} }
	_ = responseBody

	// Test
	req := httptest.NewRequest(requestMethod, requestURL, bytes.NewBuffer(requestBody))
	req.Header = requestHeader
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	if resp.StatusCode != responseStatus {
		t.Run("statuscode", func(t *testing.T) {
			t.Logf("was expecting status %q got %q instead", http.StatusText(responseStatus), http.StatusText(resp.StatusCode))
			t.FailNow()
		})
	}
	for k, _ := range resp.Header {
		t.Run(fmt.Sprintf("headers/%s", k), func(t *testing.T) {
			expected := responseHeader.Get(k)
			actual := resp.Header.Get(k)
			if expected != actual {
				t.Logf("was expecting header %q got %q instead", expected, actual)
			}
		})
	}
	if resp.Header.Get("Content-Type") == "application/json" {
		// Body only if JSON
		expected := responseBody
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Run("compareJSONresults", func(t *testing.T) {
			diffo := gojsondiff.New()
			diff, err := diffo.Compare(expected, actual)
			if err != nil {
				t.Fatal(err.Error())
			}
			if diff.Modified() {
				t.Fatal("returned json body is out of shape")
			}
		})
	}
	{{ end }}
}
	`

	t := template.Must(template.New("testTemplate").Parse(testTemplate))
	if err := t.Execute(buf, map[string]interface{}{"Trip": trip, "Name": toEnglish(trip.Request, trip.Response)}); err != nil {
		trip.Test = []byte(err.Error())
		return
	}
	opts := imports.Options{
		Fragment:   true,
		AllErrors:  false, // Report all errors (not just the first 10 on different lines)
		Comments:   true,  // Print comments (true if nil *Options provided)
		TabIndent:  true,  // Use tabs for indent (true if nil *Options provided)
		TabWidth:   8,     // Tab width (8 if nil *Options provided)
		FormatOnly: false, // Disable the insertion and deletion of imports
	}
	_ = opts
	if res, err := imports.Process("gokit.go", buf.Bytes(), &opts); err == nil {
		trip.Test = res
	} else {
		trip.Test = buf.Bytes()
	}
}

func toEnglish(req *v1.Request, resp *v1.Response) string {
	must := func(s string) *url.URL {
		u, _ := url.Parse(req.URL)
		return u
	}
	return httptestutil.ToEnglish(
		&http.Request{
			Method: req.Method,
			URL:    must(req.URL),
		},
		&http.Response{
			StatusCode: resp.Status,
		})

}

type trip struct {
	id            string
	start, finish time.Time
	req           v1.Request
	resp          v1.Response
}

var (
	specs      = make(map[string][]byte)
	bucketName = []byte("REQUESTS")
)

func (ds *debugService) saveRequest(id string, req *http.Request) {
	ds.trips[id] = trip{
		start: time.Now(),
		id:    id,
		req:   parseRequest(req),
	}
}
