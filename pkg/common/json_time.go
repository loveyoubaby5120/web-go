package common

import (
	"fmt"
	"strings"
	"time"

	"jdy/pkg/util/jsonutil"
)

const (
	// JSONTimeFormat to second precision.
	JSONTimeFormat = "2006-01-02T15:04:05"
	// JSONTimeZFormat with Z (UTC) at the end.
	JSONTimeZFormat = "2006-01-02T15:04:05Z"
	// JSONTimeMicrosFormat has microsecond precision.
	JSONTimeMicrosFormat = "2006-01-02T15:04:05.000000"
)

type JSONTime time.Time
type JSONTimeZ time.Time
type JSONTimeMicros time.Time

func JSONNow() JSONTime {
	return JSONTime(Now())
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *JSONTime) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*t = JSONTime(time.Time{})
		return nil
	}
	tt, err := time.Parse(JSONTimeFormat, string(data))
	*t = JSONTime(tt)
	return err
}

func (t JSONTime) Time() time.Time {
	return time.Time(t)
}

// ILikeDotFormat should not belong here.
func (t JSONTime) ILikeDotFormat() string {
	str := fmt.Sprintf("%4d.%2d.%2d", time.Time(t).Year(), time.Time(t).Month(), time.Time(t).Day())
	str = strings.Replace(str, " ", "0", -1)
	return str
}

func (t JSONTime) String() string {
	return timeToString(JSONTimeFormat, time.Time(t))
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return jsonutil.MarshalTime(JSONTimeFormat, time.Time(t))
}

func (t *JSONTime) UnmarshalJSON(s []byte) error {
	if len(s) == 0 {
		*t = JSONTime(time.Time{})
		return nil
	}
	if len(s) == len(JSONTimeMicrosFormat)+2 {
		tt, err := jsonutil.UnmarshalTime(JSONTimeMicrosFormat, s)
		*t = JSONTime(tt)
		return err
	}
	tt, err := jsonutil.UnmarshalTime(JSONTimeFormat, s)
	*t = JSONTime(tt)
	return err
}

// Time return time.Time
func (t JSONTimeZ) Time() time.Time {
	return time.Time(t)
}

// Before implement time.Before
func (t JSONTimeZ) Before(u JSONTimeZ) bool {
	return t.Time().Before(u.Time())
}

// Format return foramted time
func (t JSONTimeZ) Format() string {
	return t.Time().Format(JSONTimeZFormat)
}

// Between returns true if a JSONTimeZ is after startTime and before endTime
func (t JSONTimeZ) Between(startTime, endTime time.Time) bool {
	tem := t.Time()
	return tem.After(startTime) && tem.Before(endTime)
}

func (t JSONTimeZ) MarshalJSON() ([]byte, error) {
	return jsonutil.MarshalTime(JSONTimeZFormat, time.Time(t))
}

func (t JSONTimeZ) String() string {
	return timeToString(JSONTimeZFormat, time.Time(t))
}

// UnmarshalJSON offer the way to UnmarshalJSON
func (t *JSONTimeZ) UnmarshalJSON(s []byte) error {
	tt, err := jsonutil.UnmarshalTime(JSONTimeZFormat, s)
	*t = JSONTimeZ(tt)
	return err
}

func (t JSONTimeMicros) MarshalJSON() ([]byte, error) {
	return jsonutil.MarshalTime(JSONTimeMicrosFormat, time.Time(t))
}

func (t JSONTimeMicros) String() string {
	return timeToString(JSONTimeMicrosFormat, time.Time(t))
}

// UnmarshalJSON offer the way to UnmarshalJSON
func (t *JSONTimeMicros) UnmarshalJSON(s []byte) error {
	formatLen := len(s) - 2
	numMicros := formatLen - len(JSONTimeFormat)
	format := ""
	switch {
	case numMicros == 0:
		format = JSONTimeFormat
	case numMicros > 0 && numMicros <= 10:
		format = JSONTimeFormat + "." + strings.Repeat("0", numMicros-1)
	default:
		return fmt.Errorf("unknown format: %s", string(s))
	}
	tt, err := jsonutil.UnmarshalTime(format, s)
	*t = JSONTimeMicros(tt)
	return err
}

func timeToString(format string, t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(format)
}
