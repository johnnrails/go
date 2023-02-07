package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vardius/trace"
)

var (
	ErrInvalid           = errors.New("validation failed")
	ErrUnauthorized      = errors.New("access denied")
	ErrForbidden         = errors.New("forbidden")
	ErrNotFound          = errors.New("not found")
	ErrInternal          = errors.New("internal system error")
	ErrTemporaryDisabled = errors.New("temporary disabled")
	ErrTimeout           = errors.New("timeout")
)

type AppError struct {
	trace string
	err   error
}

func (e *AppError) Error() string {
	return e.err.Error()
}

func (e *AppError) Unwrap() error {
	return e.err
}

func NewAppErrorFromError(err error) *AppError {
	return newAppError(err)
}

func NewAppErrorFromMessage(m string) *AppError {
	return newAppError(errors.New(m))
}

func newAppError(err error) *AppError {
	if err == nil {
		err = ErrInternal
	}
	return &AppError{
		err:   err,
		trace: trace.FromParent(2, trace.Lfile|trace.Lline),
	}
}

// StackTrace returns the string slice of the error stack traces
func (e *AppError) StackTrace() string {
	var stack []string

	if e.trace != "" {
		stack = append(stack, e.trace)
	}

	if e.err != nil {
		var next *AppError
		if errors.As(e.err, &next) {
			stack = append(stack, next.StackTrace())
		}
	}

	return strings.Join(stack, "\n")
}

// Examples

func Example1() {
	err := NewAppErrorFromMessage("example")
	fmt.Printf("error: %s", err.Error())
	var e *AppError
	if errors.As(err, &e) {
		fmt.Printf("stack trace: %s", e.StackTrace())
	}
	// output:
	// error: example
	// stack trace: /mnt/c/Users/usuario/OneDrive/Projetos/go/ddd_go/boilerplate/http/errors/errors.go:66
}
