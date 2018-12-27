package authtest // github.com/jennyservices/jenny/auth/authtest

import (
	"context"

	"github.com/icrowley/fake"
	"github.com/jennyservices/jenny/auth"
)

type User struct {
	name, email, id string
}

func (u *User) UniqueID() []byte {
	return []byte(u.id)
}

func (u *User) Email() string {
	return u.email
}

func (u *User) DisplayName() (string, error) {
	return u.name, nil
}

func (u *User) Details() map[string]string {
	return nil
}

func FakeUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, auth.UserContextKey, &User{
		name:  fake.FullName(),
		email: fake.EmailAddress(),
		id:    fake.Characters(),
	})

}
