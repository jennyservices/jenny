package debug

import (
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Typeform/jenny/decoders"
	"github.com/Typeform/jenny/mime"
)

func init() {
	decoders.Register(mime.ApplicationOctet, newZipedJSONDecoder)
}

func newZipedJSONDecoder(r *http.Request) decoders.Decoder {
	return &zipedJSONDecoder{r: r.Body}
}

type zipedJSONDecoder struct {
	r io.Reader
}

func (z *zipedJSONDecoder) Decode(v interface{}) error {
	br := base64.NewDecoder(base64.StdEncoding, z.r)
	zr, err := zlib.NewReader(br)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(zr)
	return dec.Decode(v)
}
