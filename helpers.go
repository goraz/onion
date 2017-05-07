package onion

import (
	"reflect"
	"strings"
	"time"
)

var (
	durationValue = reflect.ValueOf(time.Second)
)

// GetStruct fill an structure base on the config nested set, this function use reflection, and its not
// good (in my opinion) for frequent call.
// but its best if you need the config to loaded in structure and use that structure after that.
func (o *Onion) GetStruct(prefix string, s interface{}) {
	iterateConfig(o, reflect.ValueOf(s), prefix)
}

func join(delim string, parts ...string) string {
	res := ""
	for i := range parts {
		if res != "" && parts[i] != "" {
			res += delim
		}
		res += parts[i]
	}

	return res
}

func setField(o *Onion, v reflect.Value, prefix, name string) {

	switch v.Kind() {
	case reflect.Bool:
		if v.CanSet() {
			v.SetBool(o.GetBoolDefault(join(o.GetDelimiter(), prefix, name), v.Bool()))
		}
	case reflect.Int:
		if v.CanSet() {
			v.SetInt(o.GetInt64Default(join(o.GetDelimiter(), prefix, name), v.Int()))
		}
	case reflect.Int64:
		if v.CanSet() {
			if v.Type().String() == durationValue.Type().String() {
				// its a duration
				v.SetInt(int64(o.GetDurationDefault(join(o.GetDelimiter(), prefix, name), time.Duration(v.Int()))))
				return
			}
			v.SetInt(o.GetInt64Default(join(o.GetDelimiter(), prefix, name), v.Int()))
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(o.GetStringDefault(join(o.GetDelimiter(), prefix, name), v.String()))
		}
	case reflect.Float64:
		if v.CanSet() {
			v.SetFloat(o.GetFloat64Default(join(o.GetDelimiter(), prefix, name), v.Float()))
		}
	case reflect.Float32:
		if v.CanSet() {
			v.SetFloat(o.GetFloat64Default(join(o.GetDelimiter(), prefix, name), v.Float()))
		}

	case reflect.Struct:
		iterateConfig(o, v.Addr(), join(o.GetDelimiter(), prefix, name))
	}
}

func iterateConfig(o *Onion, v reflect.Value, op string) {
	prefix := op
	typ := v.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		v = v.Elem()
	}
	// Only struct are supported
	if typ.Kind() != reflect.Struct {
		return
	}
	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			name := p.Tag.Get("onion")
			if name == "-" {
				// Ignore this key.
				continue
			}
			if name == "" {
				name = strings.ToLower(p.Name)
			}
			setField(o, v.Field(i), prefix, name)
		} else { // Anonymous structures
			name := p.Tag.Get("onion")
			if name == "-" {
				// Ignore this key.
				continue
			}
			if name == "" {
				prefix = op // Reset the prefix to remove the delimiter
			}
			iterateConfig(o, v.Field(i).Addr(), join(o.GetDelimiter(), prefix, name))
		}
	}

}
