# Auth

## Design

Jenny uses `go-kit` as a base framework for all the services it generates, but
it is a bit opinionated when it comes to using `go-kit` concepts.

### User

Auth package has one interface and that is User.

```golang
type User interface {
	UniqueID() []byte
}
```

It uses the `UniqueID` method for rate-limiting type operations, please make
sure the byte array that is returned from this format is unique as they will be
used to generate a unique hash. Please keep in mind that the arrays are case
sensitive so if your system of record does not differentiate between `hello` and
`Hello` first normalize the string by making all the letters uppercase or
lowercase.
