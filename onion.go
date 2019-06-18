package onion

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fzerorubigd/gozin"
)

// TODO: use cast (spf13/cast)

// Layer is an interface to handle the load phase.
type Layer interface {
	// Load is called as soon as the layer registered in the onion. if the layer is persistent
	// it can close the channel as soon as it writes the first configuration.
	// Also this function may called several time and should return the same channel each time
	// and should not block
	Load() <-chan map[string]interface{}
}

// Onion is a layer base configuration system
type Onion struct {
	lock sync.RWMutex

	delimiter string
	ll        []Layer

	// List of watches
	watches []gozin.Case
	noWatch map[Layer]bool
	stop    chan struct{}

	// Loaded data
	data map[Layer]map[string]interface{}
}

func (o *Onion) createSelectCase() []gozin.Case {
	o.lock.RLock()
	defer o.lock.RUnlock()

	var c []gozin.Case

	for i := range o.ll {
		if o.noWatch[o.ll[i]] {
			continue
		}
		l := o.ll[i]
		c = append(c,
			gozin.Receive(l.Load(), func(in interface{}, b bool) {
				if !b {
					o.closeLayerWatch(l)
					return
				}
				o.setLayerData(l, in.(map[string]interface{}))
			}))
	}

	return c
}

func (o *Onion) watchLayers(stop chan struct{}, loaded chan struct{}) {
	var (
		// done means the stop channel is closed
		done bool
	)

	var stable bool
	def := gozin.Default(func() {
		stable = true
		close(loaded)
	})
	for !done {
		c := o.createSelectCase()
		if len(c) == 0 {
			// Nothing to watch
			break
		}
		c = append(c, gozin.Receive(stop, func(_ interface{}, c bool) {
			done = !c
		}))
		if !stable {
			c = append(c, def)
		}
		if err := gozin.Select(c...); err != nil {
			panic(err)
		}
	}
	if !stable {
		close(loaded)
	}
}

func (o *Onion) setLayerData(l Layer, data map[string]interface{}) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.data == nil {
		o.data = make(map[Layer]map[string]interface{})
	}
	o.data[l] = data
}

func (o *Onion) closeLayerWatch(l Layer) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.noWatch == nil {
		o.noWatch = make(map[Layer]bool)
	}
	o.noWatch[l] = true
}

// AddLayer add a new layer to the end of config layers. last layer is loaded after all other
// layer
func (o *Onion) AddLayers(l ...Layer) {
	if len(l) == 0 {
		return
	}
	o.lock.Lock()
	o.ll = append(o.ll, l...)
	// close the old go routine
	if o.stop != nil {
		close(o.stop)
	}

	o.stop = make(chan struct{})
	o.lock.Unlock()

	loaded := make(chan struct{})
	go o.watchLayers(o.stop, loaded)
	<-loaded
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

	path := strings.Split(key, o.GetDelimiter())

	for i := len(o.ll); i > 0; i-- {
		v, ok := searchStringMap(o.data[o.ll[i-1]], path...)
		if ok {
			return v, ok
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
		// Env is not typed and always is String, so try to convert it to boolean
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

// GetStringSlice try to get a slice from the config, also it support comma separated value
// if there is no array at the key.
func (o *Onion) GetStringSlice(key string) []string {
	var ok bool
	v, ok := o.getSlice(key)
	if !ok {
		if v := o.GetString(key); len(v) > 0 {
			return strings.Split(v, ",")
		}
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
	o := &Onion{}

	o.AddLayers(layers...)
	return o
}
