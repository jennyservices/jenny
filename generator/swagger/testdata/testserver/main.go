package main

import (
	"context"
	"log"
	"net/http"

	testservice "github.com/jennyservices/jenny/testdata"
)

type fakeTesterino struct {
}

func (fakeTesterino) LoggedInUser(ctx context.Context) (*testservice.Foo, error) {
	return &testservice.Foo{}, nil
}

func main() {
	ts := testservice.NewTesterino(fakeTesterino{})
	log.Fatal(http.ListenAndServe(":8080", ts))
}
