package swagger

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorList struct {
	prefix string
	errors []error
}

func (e *ErrorList) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ErrorList) New(s string) {
	e.errors = append(e.errors, errors.Wrap(errors.New(s), e.prefix))
}
func (e *ErrorList) Add(err error) {
	e.errors = append(e.errors, errors.Wrap(err, e.prefix))
}

func (e *ErrorList) Error() string {
	s := ""
	for _, e := range e.errors {
		s += fmt.Sprintln(e)
	}
	return s
}
