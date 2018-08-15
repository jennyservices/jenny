package decoders

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/jennyservices/jenny/mime"
)

func TestFindDecoder(t *testing.T) {
	tests := []struct {
		ct      string
		body    string
		accepts []mime.Type
	}{
		{
			ct:      "application/json",
			body:    `{"hello": "world"}`,
			accepts: []mime.Type{"application/json"},
		},
	}

	for _, test := range tests {
		t.Run(test.ct, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "http://example.com", bytes.NewBufferString(test.body))
			r.Header.Add("Content-Type", test.ct)
			dec, err := RequestDecoder(r, test.accepts)
			if err != nil {
				t.Log(err)
				t.Fail()
				return
			}
			i := map[string]string{}
			if err := dec.Decode(&i); err != nil {
				t.Log(err)
				t.Fail()
			}
		})
	}
}