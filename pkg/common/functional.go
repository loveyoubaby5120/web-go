package common

// MapStringSlice string map to string slice
func MapStringSlice(s []string, f func(string) string) []string {
	if s == nil {
		return nil
	}
	o := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		o[i] = f(s[i])
	}
	return o
}

// FilterStringSlice is string filter to string slice
func FilterStringSlice(s []string, f func(string) bool) []string {
	if s == nil {
		return nil
	}
	var o []string
	for _, si := range s {
		if f(si) {
			o = append(o, si)
		}
	}
	return o
}
