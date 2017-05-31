package h_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var listToMatch = []string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"keygen",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

var result int

func BenchmarkMatchRegexp(b *testing.B) {
	patternStr := strings.Join(listToMatch, "|")
	pattern := regexp.MustCompile(fmt.Sprint("^", patternStr, "$"))
	k := 0
	for i := 0; i < b.N; i++ {
		if pattern.MatchString("div") {
			k++
		}
	}
	result = k
}

func BenchmarkMatchMap(b *testing.B) {
	m := map[string]struct{}{}
	for _, s := range listToMatch {
		m[s] = struct{}{}
	}
	k := 0
	for i := 0; i < b.N; i++ {
		if _, ok := m["div"]; ok {
			k++
		}
	}
	result = k
}

func switchM(s string) bool {
	switch s {
	case "area":
	case "base":
	case "br":
	case "col":
	case "embed":
	case "hr":
	case "img":
	case "input":
	case "keygen":
	case "link":
	case "meta":
	case "param":
	case "source":
	case "track":
	case "wbr":
		return true
	}
	return false
}

func BenchmarkMatchSwitch(b *testing.B) {
	k := 0
	for i := 0; i < b.N; i++ {
		if switchM("div") {
			k++
		}
	}
	result = k
}
