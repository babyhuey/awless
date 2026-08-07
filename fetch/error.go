package fetch

import (
	"errors"
	"strings"
)

// Not goroutine safe as for now
type Error []error

func WrapError(errs ...error) *Error {
	fe := &Error{}
	for _, e := range errs {
		var nested *Error
		if errors.As(e, &nested) {
			for _, inner := range *nested {
				fe.Add(inner)
			}
			continue
		}
		fe.Add(e)
	}
	return fe
}

func (fe *Error) Add(err error) {
	if err != nil {
		*fe = append(*fe, err)
	}
}

func (fe *Error) Any() bool {
	return len(*fe) > 0
}

func (fe *Error) Error() string {
	var all []string
	for _, e := range *fe {
		all = append(all, e.Error())
	}
	return strings.Join(all, "\n")
}
