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
	lock      sync.RWMutex
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
	if o.layers == nil {
		o.layers = make(map[Layer]map[string]interface{})
	}
	o.layers[l] = lowerStringMap(data)

	return nil
}

// GetDelimiter return the delimiter for nested key
func (o *Onion) GetDelimiter() string {
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
func (o *Onion) Get(key string) (interface{}, bool) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	if o.layers == nil {
		o.layers = make(map[Layer]map[string]interface{})
	}

	path := strings.Split(strings.ToLower(key), o.GetDelimiter())
	for i := len(o.ll) - 1; i >= 0; i-- {
		l := o.layers[o.ll[i]]
		res, found := searchStringMap(path, l)
		if found {
			return res, found
		}
	}

	return nil, false
}

// The folowing two function are identical. but converting between map[string] and
// map[interface{}] is not easy, and there is no _Generic_ , so I decide to create
// two almost identical function instead of writing a convertor each time.
//
// Some of the loaders like yaml, load inner keys in map[interface{}]interface{}
// some othr like json do it in map[string]interface{} so we should suppport both
func searchStringMap(path []string, m map[string]interface{}) (interface{}, bool) {
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch v.(type) {
	case map[string]interface{}:
		return searchStringMap(path[1:], v.(map[string]interface{}))
	case map[interface{}]interface{}:
		return searchInterfaceMap(path[1:], v.(map[interface{}]interface{}))
	}
	return nil, false
}

func searchInterfaceMap(path []string, m map[interface{}]interface{}) (interface{}, bool) {
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch v.(type) {
	case map[string]interface{}:
		return searchStringMap(path[1:], v.(map[string]interface{}))
	case map[interface{}]interface{}:
		return searchInterfaceMap(path[1:], v.(map[interface{}]interface{}))
	}
	return nil, false
}

func lowerStringMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k := range m {
		switch m[k].(type) {
		case map[string]interface{}:
			res[strings.ToLower(k)] = lowerStringMap(m[k].(map[string]interface{}))
		case map[interface{}]interface{}:
			res[strings.ToLower(k)] = lowerInterfaceMap(m[k].(map[interface{}]interface{}))
		default:
			res[strings.ToLower(k)] = m[k]
		}
	}

	return res
}

func lowerInterfaceMap(m map[interface{}]interface{}) map[interface{}]interface{} {
	res := make(map[interface{}]interface{})
	for k := range m {
		switch k.(type) {
		case string:
			switch m[k].(type) {
			case map[string]interface{}:
				res[strings.ToLower(k.(string))] = lowerStringMap(m[k].(map[string]interface{}))
			case map[interface{}]interface{}:
				res[strings.ToLower(k.(string))] = lowerInterfaceMap(m[k].(map[interface{}]interface{}))
			default:
				res[strings.ToLower(k.(string))] = m[k]
			}
		default:
			res[k] = m[k]
		}
	}

	return res
}

// GetIntDefault return an int value from Onion, if the value is not exists or its not an
// integer , default is returned
func (o *Onion) GetIntDefault(key string, def int) int {
	return int(o.GetInt64Default(key, int64(def)))
}

// GetInt return an int value, if the value is not there, then it return zero value
func (o *Onion) GetInt(key string) int {
	return o.GetIntDefault(key, 0)
}

// GetInt64Default return an int64 value from Onion, if the value is not exists or if the value is not
// int64 then return the default
func (o *Onion) GetInt64Default(key string, def int64) int64 {
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

// GetInt64 return the int64 value from config, if its not there, return zero
func (o *Onion) GetInt64(key string) int64 {
	return o.GetInt64Default(key, 0)
}

// GetStringDefault get a string from Onion. if the value is not exists or if tha value is not
// string, return the default
func (o *Onion) GetStringDefault(key string, def string) string {
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

// GetString is for getting an string from conig. if the key is not
func (o *Onion) GetString(key string) string {
	return o.GetStringDefault(key, "")
}

// GetBoolDefault return bool value from Onion. if the value is not exists or if tha value is not
// boolean, return the default
func (o *Onion) GetBoolDefault(key string, def bool) bool {
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

// GetBool is used to get a boolean value fro config, with false as default
func (o *Onion) GetBool(key string) bool {
	return o.GetBoolDefault(key, false)
}

func (o *Onion) getSlice(key string) (interface{}, bool) {
	v, ok := o.Get(key)
	if !ok {
		return nil, false
	}

	if reflect.TypeOf(v).Kind() != reflect.Slice { // Not good
		return nil, false
	}

	return v, true
}

// GetStringSlice try to get a slice from the config
func (o *Onion) GetStringSlice(key string) []string {
	var ok bool
	v, ok := o.getSlice(key)
	if !ok {
		return nil
	}

	switch v.(type) {
	case []string:
		return v.([]string)
	case []interface{}:
		vi := v.([]interface{})
		res := make([]string, len(vi))
		for i := range vi {
			if res[i], ok = vi[i].(string); !ok {
				return nil
			}
		}
		return res
	}

	return nil
}

// GetStruct fill an structure base on the config nested set, this function use reflection, and its not
// good (in my opinion) for frequent call.
// but its best if you need the config to loaded in structure and use that structure after that.
func (o *Onion) GetStruct(prefix string, s interface{}) {
	iterateConfig(o, s, prefix)
}

func iterateConfig(o *Onion, c interface{}, op string) {
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
		lock:      sync.RWMutex{},
		delimiter: ".",
		layers:    make(map[Layer]map[string]interface{}),
	}
}
