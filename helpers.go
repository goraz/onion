package onion

import (
	"reflect"
	"strings"
)

// GetStruct fill an structure base on the config nested set, this function use reflection, and its not
// good (in my opinion) for frequent call.
// but its best if you need the config to loaded in structure and use that structure after that.
func (o *Onion) GetStruct(prefix string, s interface{}) {
	iterateConfig(o, reflect.ValueOf(s), prefix)
}

func iterateConfig(o *Onion, v reflect.Value, op string) {
	prefix := op
	if prefix != "" {
		prefix = prefix + o.GetDelimiter()
	}
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

			switch v.Field(i).Kind() {
			case reflect.Bool:
				if v.Field(i).CanSet() {
					v.Field(i).SetBool(o.GetBoolDefault(prefix+name, v.Field(i).Bool()))
				}
			case reflect.Int:
				if v.Field(i).CanSet() {
					v.Field(i).SetInt(o.GetInt64Default(prefix+name, v.Field(i).Int()))
				}
			case reflect.Int64:
				if v.Field(i).CanSet() {
					v.Field(i).SetInt(o.GetInt64Default(prefix+name, v.Field(i).Int()))
				}
			case reflect.String:
				if v.Field(i).CanSet() {
					v.Field(i).SetString(o.GetStringDefault(prefix+name, v.Field(i).String()))
				}
			case reflect.Float64:
				if v.Field(i).CanSet() {
					v.Field(i).SetFloat(o.GetFloat64Default(prefix+name, v.Field(i).Float()))
				}
			case reflect.Float32:
				if v.Field(i).CanSet() {
					v.Field(i).SetFloat(o.GetFloat64Default(prefix+name, v.Field(i).Float()))
				}
			case reflect.Struct:
				iterateConfig(o, v.Field(i).Addr(), prefix+name)
			}
		} else { // Anonymous structures
			name := p.Tag.Get("onion")
			if name == "" {
				prefix = op // Reset the prefix to remove the delimiter
			}
			iterateConfig(o, v.Field(i).Addr(), prefix+name)
		}
	}

}
