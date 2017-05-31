package must

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// String returns the string or panic.
func String(s string, err error) string {
	Must(err)
	return s
}

// NotEmpty checks string not empty.
func NotEmpty(s string) {
	if s == "" {
		panic("given string is empty")
	}
}

// True checks b is true.
func True(b bool) {
	if !b {
		panic("assertion not true")
	}
}
