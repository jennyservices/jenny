---
id: user-service-2
title: User Service-2
sidebar_label: User Service Tutorial-2 Adding a datastore
---

## Getting some use out of our user package

Let's start by adding a datastore package to our service by creating
`datastore/datastore.go`

In this package we'll define 2 things;

1. `Users` datastore interface
2. `ErrUserNotFound` error

```go
package datastore

import (
  "github.com/Typeform/users/transport/v1"
  "errors"
)

var (
  // ErrUserNotFound is returned when a DB can't find a user
  ErrUserNotFound = errors.New("user not found")
)

// Users represents a backing datastore for the Users service
type Users interface {
	GetUser(int) (*v1.User, error)
}
```

We'll now switch over to the `Users` implementation and make sure it's `New`
function expects a `datastore.Users` implementation.

```go
type userService struct {
	db datastore.Users
}

// New returns a new v1.Users as required by the transport layer, it accepts
// a datastore.Users as it's only argument.
func New(db datastore.Users) v1.Users {
	return &userService{
		db: db,
	}
}

func (us *userService) GetUser(ctx context.Context, id int) (*v1.User, error) {
	return us.db.GetUser(id)
}
```

Now let's go ahead and implement users datastore in memory

```go
package memstore

import (
	"github.com/Typeform/users/datastore"
	"github.com/Typeform/users/transport/v1"
)

// Memdb is a in-memory users.DB implementation
type memdb struct {
	users map[int]v1.User
}

func (m *memdb) GetUser(id int) (*v1.User, error) {
	if user, ok := m.users[id]; ok {
		return &user, nil
	}
	return nil, datastore.ErrUserNotFound
}

// New returns an in-memory datastore.Users
func New(s map[int]v1.User) datastore.Users {
	return &memdb{
		users: s,
	}
}
```

And let's glue everything togheter

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Typeform/users/datastore/memstore"
	"github.com/Typeform/users/transport/v1"
	"github.com/Typeform/users/user"
)

var (
	addr = flag.String("addr", ":8080", "default -addr :8080")
)

func Int(i int) *int {
	return &i
}

func String(s string) *string {
	return &s
}

func main() {
	flag.Parse()

	x := map[int]v1.User{
		1: v1.User{
			ID:   Int(1),
			Name: String("sevki"),
		},
	}
	userService := user.New(memstore.New(x))
	us := v1.NewUsersHTTPServer(userService)

	log.Printf("listening at %s\n", *addr)
	log.Fatal("listening at "+*addr+" failed: ", http.ListenAndServe(*addr, us))
}
```

[See all files](https://github.com/sevki/user/tree/tutorial-2)
