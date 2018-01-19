package debug

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/Typeform/jenny/errors"

	"github.com/Typeform/jenny/debug/transport/v1"
	"github.com/boltdb/bolt"
)

const errbktname = "err"

func makeID(i uint64, id string) []byte {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s_%d", id, i)
	return buf.Bytes()
}

func (ds *debugService) NewError(ctx context.Context, ID string, Packet v1.Packet) (*v1.ErrorResponse, error) {
	if Packet.Extra.XRequestID == nil {
		return nil, errors.NewHTTPError(errors.New("no request id"), http.StatusInternalServerError)
	}
	err := ds.db.Update(func(tx *bolt.Tx) error {
		errorBucket, err := tx.CreateBucketIfNotExists([]byte(errbktname))
		if err != nil {
			log.Println(err)
			return err
		}

		buf := bytes.NewBuffer(nil)
		enc := gob.NewEncoder(buf)
		if err := enc.Encode(Packet); err != nil {
			log.Println(err)
			return err
		}
		id, _ := errorBucket.NextSequence()

		return errorBucket.Put(makeID(id, *Packet.Extra.XRequestID), buf.Bytes())
	})
	if err != nil {
		log.Fatal(err)
	}
	return &v1.ErrorResponse{Wrote: true}, err
}
