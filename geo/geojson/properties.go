package geojson

// Properties defines the feature properties with some helper methods.
type Properties map[string]interface{}

// MustBool guarantees the return of a `bool` (with optional default).
// This function useful when you explicitly want a `bool` in a single
// value return context, for example:
//     myFunc(f.Properties.MustBool("param1"), f.Properties.MustBool("optional_param", true))
func (p Properties) MustBool(key string, def ...bool) bool {
	b, ok := p[key].(bool)
	if ok {
		return b
	}

	if len(def) > 0 {
		return def[0]
	}

	return false
}

// MustInt guarantees the return of an `int` (with optional default).
// This function useful when you explicitly want a `int` in a single
// value return context, for example:
//     myFunc(f.Properties.MustInt("param1"), f.Properties.MustInt("optional_param", 123))
func (p Properties) MustInt(key string, def ...int) int {
	i, ok := p[key].(int)
	if ok {
		return i
	}

	f, ok := p[key].(float64)
	if ok {
		return int(f)
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

// MustFloat64 guarantees the return of a `float64` (with optional default)
// This function useful when you explicitly want a `float64` in a single
// value return context, for example:
//     myFunc(f.Properties.MustFloat64("param1"), f.Properties.MustFloat64("optional_param", 10.1))
func (p Properties) MustFloat64(key string, def ...float64) float64 {
	f, ok := p[key].(float64)
	if ok {
		return f
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0.0
}

// MustString guarantees the return of a `string` (with optional default)
// This function useful when you explicitly want a `string` in a single
// value return context, for example:
//     myFunc(f.Properties.MustString("param1"), f.Properties.MustString("optional_param", "default"))
func (p Properties) MustString(key string, def ...string) string {
	s, ok := p[key].(string)
	if ok {
		return s
	}

	if len(def) > 0 {
		return def[0]
	}

	return ""
}
