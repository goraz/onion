package onion

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Layer is an interface to handle the load phase.
type Layer interface {
	// Reload return true if there is a change in the data
	Reload() (bool, error)

	// Get should return the key, each call is only for one key. for example Get(foo, bar) means the key
	// bar inside the root key foo
	Get(...string) (interface{}, bool)
}

type layerList []Layer

// Onion is a layer base configuration system
type Onion struct {
	lock      sync.RWMutex
	delimiter string
	ll        layerList

	// Cache for a layer, should change after reload
	cache map[string]interface{}
}

func (o *Onion) setCache(key string, val interface{}) {
	if o.cache == nil {
		o.cache = make(map[string]interface{})
	}

	o.cache[key] = val
}

func (o *Onion) getCache(key string) (interface{}, bool) {
	v, ok := o.cache[key]
	return v, ok
}

// AddLayer add a new layer to the end of config layers. last layer is loaded after all other
// layer
func (o *Onion) AddLayer(l Layer) error {
	o.lock.Lock()
	defer o.lock.Unlock()

	o.ll = append(o.ll, l)
	// reset the cache
	o.cache = make(map[string]interface{})

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

	if v, b := o.getCache(key); b {
		return v, true
	}

	path := strings.Split(key, o.GetDelimiter())
	for i := len(o.ll); i > 0; i-- {
		if v, b := o.ll[i-1].Get(path...); b {
			o.setCache(key, v)
			return v, true
		}
	}

	return nil, false
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
		return nil
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
func New(layers ...Layer) *Onion {
	return &Onion{
		lock:      sync.RWMutex{},
		delimiter: ".",
		cache:     make(map[string]interface{}),
		ll:        layers,
	}
}
