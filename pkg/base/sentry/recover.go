package sentry

import "fmt"

// Recover recovers panic and report it to sentry. And re-panic.
// Must be called directly by a defer, like this: defer sentry.Recover(false)
func Recover(repanic bool) {
	if r := recover(); r != nil {
		Error(fmt.Errorf("Recover panic: %v", r))
		if repanic {
			panic(r)
		}
	}
}
