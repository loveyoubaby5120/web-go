package jsonutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MarshalTime marshal a time into json.
func MarshalTime(format string, t time.Time) ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Format(format))), nil
}

// UnmarshalTime unmarshal a raw json to time with given format.
func UnmarshalTime(format string, s []byte) (time.Time, error) {
	if string(s) == "null" || string(s) == `""` {
		return time.Time{}, nil
	}
	q, err := strconv.Unquote(string(s))
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(format, q, time.UTC)
}

// UnmarshalTimeMultiFormat tries to parse the given json in different format.
func UnmarshalTimeMultiFormat(formats []string, s []byte) (time.Time, error) {
	if string(s) == "null" || string(s) == `""` {
		return time.Time{}, nil
	}
	q, err := strconv.Unquote(string(s))
	if err != nil {
		return time.Time{}, err
	}

	for _, format := range formats {
		t, err := time.ParseInLocation(format, q, time.UTC)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf(
		"%s doesn't match any of the time formats: %s",
		string(s), strings.Join(formats, ", "))
}
