package common

import (
	"fmt"
	"regexp"
	"strings"
)

func SplitAndTrim(s string, d string) []string {
	return FilterStringSlice(
		MapStringSlice(strings.Split(s, d), strings.TrimSpace),
		func(s string) bool {
			return s != ""
		})
}

func FmtMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		if v, ok := msgAndArgs[0].(string); ok {
			return v
		}
		if v, ok := msgAndArgs[0].(error); ok {
			return v.Error()
		}
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

func NewStrPtr(s string) *string {
	return &s
}

// EmailPattern email pattern
var EmailPattern = regexp.MustCompile(`^([a-zA-Z0-9_\.\-])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})$`)

// ValidEmail if the given email is valid, return true
func ValidEmail(email string) bool {
	return EmailPattern.MatchString(email)
}

// GetSourceByEmail figure out which source a email belongs to
func GetSourceByEmail(email string) string {
	if ValidEmail(email) {
		if strings.HasSuffix(email, "@trial.applysquare.com") {
			return "visitor"
		}
		switch strings.Split(strings.Split(email, "@")[1], ".")[0] {
		case "qq":
			return "qq"
		case "weibo":
			return "weibo"
		case "wechat":
			return "wechat"
		default:
			return "email"
		}
	}
	return "Not an email"
}
