package i18n

//go:generate go run ../../cmd/manong/manong.go --in_dir=.. --translation_file_path=zh.go,en.go --logtostderr collect_translation

import (
	"fmt"
	"strings"
)

type Language int

const EnUS Language = 0
const ZhCN Language = 1

var EnContext = &Context{EnUS}
var CnContext = &Context{ZhCN}

var stringMap = func() []map[string]string {
	m := make([]map[string]string, 2)
	m[EnUS] = GetENMap()
	m[ZhCN] = GetZHMap()
	return m
}()

type Context struct {
	lang Language
}

func NewContext(language string) *Context {
	switch strings.ToLower(language) {
	case "zh-cn", "cn", "zh":
		return CnContext
	}
	return EnContext
}

func (c *Context) Language() Language {
	return c.lang
}

func (c *Context) LangCode() string {
	switch c.lang {
	case ZhCN:
		return "zh-cn"
	default:
		return "en-us"
	}
}

func (c *Context) ShortLangCode() string {
	switch c.lang {
	case ZhCN:
		return "cn"
	default:
		return "en"
	}
}

func (c *Context) S(s string, a ...interface{}) string {
	translated := func() string {
		if c.lang > 1 || c.lang < 0 {
			return s
		}
		m := stringMap[c.lang]
		r := m[s]
		if r != "" {
			return r
		}
		r = m[toLower(s)]
		if r != "" {
			return r
		}
		return s
	}()
	if len(a) == 0 {
		return translated
	}
	return fmt.Sprintf(translated, a...)
}

// TODO: Replace S with SS.

func (c *Context) SS(s string, a ...interface{}) S {
	return S(c.S(s, a...))
}

func (c *Context) ConcatS(as ...interface{}) S {
	return S(fmt.Sprint(as...))
}

// ToTrans is a marker to mark the string should be picked up
// for translation. We typically shouldn't need it.
func ToTrans(s string) string {
	return s
}

// toLower is much faster than strings.ToLower.
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
