package onion

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultDelimiter is the default delimiter for the config scope
const DefaultDelimiter = "."

var lock = &sync.RWMutex{}

// Layer is an interface to handle the load phase.
type Layer interface {
	// Load a layer into the Onion. the call is only done in the
	// registration
	Load() (map[string]interface{}, error)
}

// LazyLayer is the layer for lazy config sources, when the entire configs is not
// available at the registration
type LazyLayer interface {
	// Get return the value for this config in this layer, if exists, if not return
	// false as the 2nd return value
	Get(...string) (interface{}, bool)
}

type singleLayer struct {
	layer Layer
	lazy  LazyLayer
	data  map[string]interface{}
}

type layerList []singleLayer

// Onion is a layer base configuration system
type Onion struct {
	delimiter string
	ll        layerList
}

func (sl singleLayer) getData(d string, path ...string) (interface{}, bool) {
	if sl.lazy != nil {
		return sl.lazy.Get(path...)
	}
	return searchStringMap(path, sl.data)
}

// AddLayer add a new layer to the end of config layers. last layer is loaded after all other
// layer
func (o *Onion) AddLayer(l Layer) error {
	lock.Lock()
	defer lock.Unlock()

	data, err := l.Load()
	if err != nil {
		return err
	}

	o.ll = append(
		o.ll,
		singleLayer{
			layer: l,
			lazy:  nil,
			data:  lowerStringMap(data),
		},
	)

	return nil
}

// AddLazyLayer add a new lazy layer to the end of config layers. last layer is loaded after
// all other layer
func (o *Onion) AddLazyLayer(l LazyLayer) {
	lock.Lock()
	defer lock.Unlock()

	o.ll = append(
		o.ll,
		singleLayer{
			layer: nil,
			lazy:  l,
			data:  nil,
		},
	)
}

// GetDelimiter return the delimiter for nested key
func (o *Onion) GetDelimiter() string {
	if o.delimiter == "" {
		o.delimiter = DefaultDelimiter
	}

	return o.delimiter
}

// SetDelimiter set the current delimiter
func (o *Onion) SetDelimiter(d string) {
	o.delimiter = d
}

// Get try to get the key from config layers
func (o *Onion) Get(key string) (interface{}, bool) {
	lock.RLock()
	defer lock.RUnlock()
	key = strings.Trim(key, " ")
	if len(key) == 0 {
		return nil, false
	}
	path := strings.Split(strings.ToLower(key), o.GetDelimiter())
	for i := len(o.ll) - 1; i >= 0; i-- {
		res, found := o.ll[i].getData(o.GetDelimiter(), path...)
		if found {
			return res, found
		}
	}

	return nil, false
}

// The following two function are identical. but converting between map[string] and
// map[interface{}] is not easy, and there is no _Generic_ way to do it, so I decide to create
// two almost identical function instead of writing a converter each time.
//
// Some of the loaders like yaml, load inner keys in map[interface{}]interface{}
// some other like json do it in map[string]interface{} so we should support both
func searchStringMap(path []string, m map[string]interface{}) (interface{}, bool) {
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(path[1:], m)
	case map[interface{}]interface{}:
		return searchInterfaceMap(path[1:], m)
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

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(path[1:], m)
	case map[interface{}]interface{}:
		return searchInterfaceMap(path[1:], m)
	}
	return nil, false
}

func lowerStringMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k := range m {
		switch nm := m[k].(type) {
		case map[string]interface{}:
			res[strings.ToLower(k)] = lowerStringMap(nm)
		case map[interface{}]interface{}:
			res[strings.ToLower(k)] = lowerInterfaceMap(nm)
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
			switch nm := m[k].(type) {
			case map[string]interface{}:
				res[strings.ToLower(k.(string))] = lowerStringMap(nm)
			case map[interface{}]interface{}:
				res[strings.ToLower(k.(string))] = lowerInterfaceMap(nm)
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

	switch nv := v.(type) {
	case string:
		// Env is not typed and always is String, so try to convert it to int
		// if possible
		i, err := strconv.ParseInt(nv, 10, 64)
		if err != nil {
			return def
		}
		return i
	case int:
		return int64(nv)
	case int64:
		return nv
	case float32:
		return int64(nv)
	case float64:
		return int64(nv)
	default:
		return def
	}

}

// GetInt64 return the int64 value from config, if its not there, return zero
func (o *Onion) GetInt64(key string) int64 {
	return o.GetInt64Default(key, 0)
}

// GetFloat32Default return an float32 value from Onion, if the value is not exists or its not a
// float32, default is returned
func (o *Onion) GetFloat32Default(key string, def float32) float32 {
	return float32(o.GetFloat64Default(key, float64(def)))
}

// GetFloat32 return an float32 value, if the value is not there, then it returns zero value
func (o *Onion) GetFloat32(key string) float32 {
	return o.GetFloat32Default(key, 0)
}

// GetFloat64Default return an float64 value from Onion, if the value is not exists or if the value is not
// float64 then return the default
func (o *Onion) GetFloat64Default(key string, def float64) float64 {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	switch nv := v.(type) {
	case string:
		// Env is not typed and always is String, so try to convert it to int
		// if possible
		f, err := strconv.ParseFloat(nv, 64)
		if err != nil {
			return def
		}
		return f
	case int:
		return float64(nv)
	case int64:
		return float64(nv)
	case float32:
		return float64(nv)
	case float64:
		return nv
	default:
		return def
	}

}

// GetFloat64 return the float64 value from config, if its not there, return zero
func (o *Onion) GetFloat64(key string) float64 {
	return o.GetFloat64Default(key, 0)
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

	switch nv := v.(type) {
	case string:
		// Env is not typed and always is String, so try to convert it to int
		// if possible
		i, err := strconv.ParseBool(nv)
		if err != nil {
			return def
		}
		return i
	case bool:
		return nv
	default:
		return def
	}
}

// GetBool is used to get a boolean value fro config, with false as default
func (o *Onion) GetBool(key string) bool {
	return o.GetBoolDefault(key, false)
}

// GetDurationDefault is a function to get duration from config. it support both
// string duration (like 1h3m2s) and integer duration
func (o *Onion) GetDurationDefault(key string, def time.Duration) time.Duration {
	v, ok := o.Get(key)
	if !ok {
		return def
	}

	switch nv := v.(type) {
	case string:
		d, err := time.ParseDuration(nv)
		if err != nil {
			return def
		}
		return d
	case int:
		return time.Duration(nv)
	case int64:
		return time.Duration(nv)
	case time.Duration:
		return nv
	default:
		return def
	}
}

// GetDuration is for getting duration from config, it cast both int and string
// to duration
func (o *Onion) GetDuration(key string) time.Duration {
	return o.GetDurationDefault(key, 0)
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
		str := o.GetString(key)
		if len(str) <= 0 {
			return nil
		}
		v = strings.Split(str, ",")
	}

	switch nv := v.(type) {
	case []string:
		return nv
	case []interface{}:
		res := make([]string, len(nv))
		for i := range nv {
			if res[i], ok = nv[i].(string); !ok {
				return nil
			}
		}
		return res
	}

	return nil
}

// New return a new Onion
func New() *Onion {
	return &Onion{
		delimiter: DefaultDelimiter,
	}
}
