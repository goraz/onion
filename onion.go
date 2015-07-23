package onion

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Layer is an interface to handle the load phase.
type Layer interface {
	// Load a layer into the Onion
	Load() (map[string]interface{}, error)
}

type layerList []Layer

// Onion is a layer base configuration system
type Onion struct {
	lock      *sync.Mutex
	delimiter string
	ll        layerList
	// A simple cache system, should have a way to refresh
	layers map[Layer]map[string]interface{}
}

// AddLayer add a new layer to the end of config layers. last layer is loaded after all other
// layer
func (o *Onion) AddLayer(l Layer) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	data, err := l.Load()
	if err != nil {
		return err
	}

	o.ll = append(o.ll, l)
	o.layers[l] = lowerMap(data)

	return nil
}

// GetDelimiter return the delimiter for nested key
func (o Onion) GetDelimiter() string {
	if o.delimiter == "" {
		o.delimiter = "."
	}

	return o.delimiter
}

// SetDelimiter set the current delimiter
func (o *Onion) SetDelimiter(d string) {
	o.delimiter = d
}

// Get try to get the key from config layers
func (o Onion) Get(key string) (interface{}, bool) {
	path := strings.Split(strings.ToLower(key), o.GetDelimiter())
	for i := len(o.ll) - 1; i >= 0; i-- {
		l := o.layers[o.ll[i]]
		res, found := searchMap(path, l)
		if found {
			return res, found
		}
	}

	return nil, false
}

func searchMap(path []string, m map[string]interface{}) (interface{}, bool) {
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch v.(type) {
	case map[string]interface{}:
		return searchMap(path[1:], v.(map[string]interface{}))
	}
	return nil, false
}

func lowerMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k := range m {
		switch m[k].(type) {
		case map[string]interface{}:
			res[strings.ToLower(k)] = lowerMap(m[k].(map[string]interface{}))
		default:
			res[strings.ToLower(k)] = m[k]
		}
	}

	return res
}

// GetInt return an int value from Onion, if the value is not exists or its not an
// integer , default is returned
func (o Onion) GetInt(key string, def int) int {
	return int(o.GetInt64(key, int64(def)))
}

// GetInt64 return an int64 value from Onion, if the value is not exists or if the value is not
// int64 then return the default
func (o Onion) GetInt64(key string, def int64) int64 {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	switch v.(type) {
	case string:
		// Env is not typed and always is String, so try to convert it to int
		// if possible
		i, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return def
		}
		return i
	case int:
		return int64(v.(int))
	case int64:
		return v.(int64)
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	default:
		return def
	}
}

// GetString get a string from Onion. if the value is not exists or if tha value is not
// string, return the default
func (o Onion) GetString(key string, def string) string {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	s, ok := v.(string)
	if !ok {
		return def
	}

	return s
}

// GetBool return bool value from Onion. if the value is not exists or if tha value is not
// boolean, return the default
func (o Onion) GetBool(key string, def bool) bool {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	switch v.(type) {
	case string:
		// Env is not typed and always is String, so try to convert it to int
		// if possible
		i, err := strconv.ParseBool(v.(string))
		if err != nil {
			return def
		}
		return i
	case bool:
		return v.(bool)
	default:
		return def
	}
}

// GetStruct fill an structure base on the config nested set
func (o Onion) GetStruct(s interface{}) {
	iterateConfig(o, s, "")
}

func iterateConfig(o Onion, c interface{}, op string) {
	prefix := op
	if prefix != "" {
		prefix = prefix + o.GetDelimiter()
	}
	typ := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		v = v.Elem()
	}
	// Only structs are supported
	if typ.Kind() != reflect.Struct {
		return
	}

	// loop through the struct's fields and set the map
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			name := p.Tag.Get("onion")
			if name == "" {
				name = strings.ToLower(p.Name)
			}

			switch v.Field(i).Kind() {
			case reflect.Bool:
				if v.Field(i).CanSet() {
					v.Field(i).SetBool(o.GetBool(prefix+name, v.Field(i).Bool()))
				}
			case reflect.Int:
				if v.Field(i).CanSet() {
					v.Field(i).SetInt(o.GetInt64(prefix+name, v.Field(i).Int()))
				}
			case reflect.Int64:
				if v.Field(i).CanSet() {
					v.Field(i).SetInt(o.GetInt64(prefix+name, v.Field(i).Int()))
				}
			case reflect.String:
				if v.Field(i).CanSet() {
					v.Field(i).SetString(o.GetString(prefix+name, v.Field(i).String()))
				}
			case reflect.Struct:
				iterateConfig(o, v.Field(i).Addr().Interface(), prefix+name)
			}
		} else { // Anonymus structues
			name := p.Tag.Get("onion")
			if name == "" {
				prefix = op // Reset the prefix to remove the delimiter
			}
			iterateConfig(o, v.Field(i).Addr().Interface(), prefix+name)
		}
	}

}

// New return a new Onion
func New() *Onion {
	return &Onion{
		lock:      &sync.Mutex{},
		delimiter: ".",
		layers:    make(map[Layer]map[string]interface{}),
	}
}
