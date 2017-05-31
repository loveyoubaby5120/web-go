package common

import (
	"sync/atomic"
	"time"

	"jdy/pkg/util/errs"
	"jdy/pkg/util/must"
)

var (
	fakeTime atomic.Value
)

func init() {
	fakeTime.Store(time.Time{})
}

func SetFakeTime(t time.Time) {
	fakeTime.Store(t)
}

func AdvanceFakeTime(d time.Duration) {
	t := fakeTime.Load().(time.Time)
	t = t.Add(d)
	fakeTime.Store(t)
}

func SetAFakeTime() {
	t, e := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41.123456Z")
	must.Must(e)
	SetFakeTime(t)
}

func UnSetFakeTime() {
	fakeTime.Store(time.Time{})
}

func Now() time.Time {
	t := fakeTime.Load().(time.Time)
	if !t.IsZero() {
		return t
	}
	return time.Now().UTC()
}

func ParseTimeOrInvalidError(s string, format string) (time.Time, error) {
	value, err := time.Parse(format, s)
	if err != nil {
		return value, errs.InvalidArgument.Wrapf(err, "Failed to parse time-string %s against %s", s, format)
	}
	return value, nil
}
