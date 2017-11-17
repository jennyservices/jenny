// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"crypto/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

func getID(r *http.Request) string {
	reqID := r.Header.Get("X-Request-Id")
	if reqID != "" {
		return reqID
	}
	id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}

func (ds *debugService) debugHandler(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		id := getID(r)
		r.Header.Set("X-Debug-ID", id)
		ds.saveRequest(id, r)
		old := w
		w = NewRecorder(w)
		h.ServeHTTP(w, r)
		resp := w.(*ResponseRecorder).Result()
		ds.saveResponse(id, resp)
		for k, vs := range w.(*ResponseRecorder).result.Header {
			for _, v := range vs {
				old.Header().Add(k, v)
			}
		}
		old.WriteHeader(resp.StatusCode)
		old.Write(w.(*ResponseRecorder).Body.Bytes())
		return
	}
	return http.HandlerFunc(f)
}
