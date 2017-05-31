package sentry_test

import (
	"testing"

	"jdy/pkg/base/sentry"
)

func TestRecover(t *testing.T) {
	func() {
		defer sentry.Recover(false)
		panic("hi")
	}()
}
