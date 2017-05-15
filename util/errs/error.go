package errs

import (
	"bytes"
	"fmt"
)

// Error implements error interface and capture the relevant information.
type Error struct {
	Kind    Kind
	Message string
	Err     error // Underlying error.

	stack stack
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	printStack(e.stack, b)

	pad(b, "\n")
	fmt.Fprintf(b, "%s errror", e.Kind.String())
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	b.WriteString(".")
	if e.Err != nil {
		pad(b, "\n")
		fmt.Fprintf(b, "Caused by: %v", e.Err)
	}
	return b.String()
}

// Wrap wraps the error.
func Wrap(err error) error {
	return WrapfSkipFrame(1, err, "")
}

// Wrapf is almost equivalent to Internal.Wrapf, but it tries to
// preserve the kind, and ignores nil error.
func Wrapf(err error, msg string, args ...interface{}) error {
	return WrapfSkipFrame(1, err, msg, args...)
}

// WrapfSkipFrame is almost equivalent to Internal.Wrapf, but it tries to
// preserve the kind, and ignores nil error.
func WrapfSkipFrame(depth int, err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if e2, ok := err.(*Error); ok {
		return e2.Kind.WrapfSkipFrame(1+depth, err, msg, args...)
	}
	return Internal.WrapfSkipFrame(1+depth, err, msg, args...)
}

// New is equivalent to Internal.New.
func New(msg string, args ...interface{}) error {
	return Internal.WrapfSkipFrame(1, nil, msg, args...)
}

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ":\n\t"

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
