package errs

import "fmt"
import "net/http"

// Kind defines the kind of error.
type Kind uint8

const (
	// Internal is the generic error that maps to HTTP 500.
	Internal Kind = iota
	// NotFound indicates a given resource is not found.
	NotFound
	// Forbidden indicates the user doesn't have the permission to
	// perform given operation.
	Forbidden
	// Unauthenticated indicates the user hasn't login.
	Unauthenticated
	// InvalidArgument indicates the input is invalid.
	InvalidArgument
	// Conflict indicates a database transactional conflict happens.
	Conflict
	// TryAgain indicates a temporary outage and retry
	// could eventually lead to success.
	TryAgain
)

// String retruns name of the kind.
func (k Kind) String() string {
	switch k {
	case Internal:
		return "Internal"
	case NotFound:
		return "NotFound"
	case Forbidden:
		return "Forbidden"
	case Unauthenticated:
		return "Unauthenticated"
	case InvalidArgument:
		return "InvalidArgument"
	case Conflict:
		return "Conflict"
	case TryAgain:
		return "TryAgain"
	}
	panic(fmt.Errorf("unknown kind: %v", int(k)))
}

// Wrap wraps a error with given kind.
func (k Kind) Wrap(err error) error {
	return k.WrapfSkipFrame(1, err, "")
}

// Wrapf wraps a error with given kind.
func (k Kind) Wrapf(err error, msg string, args ...interface{}) error {
	return k.WrapfSkipFrame(1, err, msg, args...)
}

// WrapfSkipFrame wraps with extra skip frames.
func (k Kind) WrapfSkipFrame(depth int, err error, msg string, args ...interface{}) error {
	if msg == "" && len(args) == 0 {
		if e, ok := err.(*Error); ok && e.Kind == k {
			return err
		}
	}
	e := &Error{
		Kind:    k,
		Message: fmt.Sprintf(msg, args...),
		Err:     err,
		stack:   callers(depth),
	}
	mergeStack(e)
	return e
}

// New is equivalent to Wrap without underlying error.
func (k Kind) New(msg string, args ...interface{}) error {
	return k.WrapfSkipFrame(1, nil, msg, args...)
}

// Is returns whether given error has the given type.
func (k Kind) Is(err error) bool {
	errV, ok := err.(*Error)
	if !ok {
		// All other errors are internal error.
		return k == Internal
	}
	return errV.Kind == k
}

// HTTPStatusCode returns the HTTP status code of given kind.
func (k Kind) HTTPStatusCode() int {
	switch k {
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case Forbidden:
		return http.StatusForbidden
	case Unauthenticated:
		return http.StatusUnauthorized
	case InvalidArgument:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case TryAgain:
		return http.StatusServiceUnavailable
	default:
		panic(fmt.Errorf("unknown kind: %v", k))
	}
}
