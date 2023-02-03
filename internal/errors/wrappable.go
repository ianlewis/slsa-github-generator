package errors

import (
	stderrors "errors"
	"fmt"
)

var (
	// New is the same as errors.New.
	New = stderrors.New

	// As is the same as errors.As.
	As = stderrors.As
)

// Wrappable is a wrappable error.
type Wrappable interface {
	Error() string
	setWrapped(error)
}

// WrappableError is a wrappable struct that can be easily embedded in error
// types.
type WrappableError struct {
	err error
}

func (e *WrappableError) Error() string {
	wrapped := e.Unwrap()
	if wrapped == nil {
		return ""
	}
	return wrapped.Error()
}

func (e *WrappableError) Unwrap() error {
	return e.err
}

func (e *WrappableError) setWrapped(err error) {
	e.err = err
}

// Errorf returns a wrappable typed error that can be type checked in tests.
func Errorf(err Wrappable, format string, a ...any) error {
	err.setWrapped(fmt.Errorf(format, a...))
	return err
}

// CheckAs wraps As allowing the caller to pass an error instance rather than a
// concrete error type. This is useful for table-driven tests where errors can
// be stored in error instances but checked with their concrete types.
func CheckAs[T error](got error, want T) bool {
	return As(got, &want)
}
