package debug

import (
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/Typeform/jenny/decoders"
	"github.com/Typeform/jenny/mime"
)

func init() {
	decoders.Register(mime.ApplicationOctet, newZipedJSONDecoder)
}

func newZipedJSONDecoder(r io.Reader) decoders.Decoder {
	return &zipedJSONDecoder{r: r}
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
