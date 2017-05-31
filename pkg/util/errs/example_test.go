package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"jdy/pkg/util/errs"
)

var realErr = errors.New("real")

func makeError() error {
	return errs.InvalidArgument.Wrapf(nil, "hi")
}

func TestErrors(t *testing.T) {
	err := makeError()
	assert.Error(t, err)

	err = errs.Wrap(err)
	assert.Error(t, err)
	assert.True(t, errs.InvalidArgument.Is(err))

	err = errs.NotFound.New("not found")
	assert.True(t, errs.NotFound.Is(err))
	assert.False(t, errs.Forbidden.Is(err))

	err = errs.Wrap(realErr)
	assert.True(t, err != realErr)
	assert.True(t, errs.Unwrap(err) == realErr)

	err = errs.InvalidArgument.Wrap(err)
	assert.True(t, errs.Unwrap(err) == realErr)
	assert.True(t, errs.InvalidArgument.Is(err))
}
