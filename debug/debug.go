// Copyright 2017 Typeform SL. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"

	"github.com/Typeform/jenny/debug/transport/v1"
	"github.com/boltdb/bolt"
)

func (ds *debugService) Trips(ctx context.Context) ([]v1.Trip, error) {
	trps := []v1.Trip{}
	ds.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucketName)
		b.ForEach(func(k, v []byte) error {
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)

			t := v1.Trip{}
			dec.Decode(&t)
			t.Request = nil
			t.Response = nil
			//createTest(t)
			trps = append(trps, t)
			return nil
		})
		return nil
	})

	sort.Slice(trps, func(i int, j int) bool {
		return ds.trips[trps[i].ID].start.After(ds.trips[trps[j].ID].start)
	})
	return trps, nil
}

func (ds *debugService) GetTrip(ctx context.Context, id string) (*v1.Trip, error) {
	t := new(v1.Trip)
	err := ds.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucketName)
		buf := bytes.NewBuffer(b.Get([]byte(id)))
		dec := gob.NewDecoder(buf)
		dec.Decode(t)
		createTest(t)
		return nil
	})

	return t, err
}

func urlToString(u *url.URL) *string {
	if u != nil {
		s := u.String()
		return &s
	}
	return nil
}

func (ds *debugService) Run(ctx context.Context) (*v1.Run, error) {
	if ds.proxy == nil {
		return nil, errors.New("can't find the builder")
	}

	return &v1.Run{
		Builder: ds.proxy.BuildLog(),
		Runner:  ds.proxy.RunLog(),
		Addr:    urlToString(ds.proxy.Listens()),
	}, nil
}

func (ds *debugService) GetService(ctx context.Context, name string) ([]byte, error) {
	resp, err := http.Get("http://0.0.0.0:8081/_swagger")
	if err != nil {
		return nil, err
	}
	bytz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return bytz, nil
}

func (ds *debugService) Services(ctx context.Context) ([]v1.Service, error) {
	services := []v1.Service{}

	return services, nil
}
