package onion

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Layer is an interface to handle the load phase.
type Layer interface {
	// Load is called once to get the initial data, it can return nil if there is no initial data
	Load() map[string]interface{}
	// Watch is called as soon as the layer registered in the onion. if the layer is persistent
	// it can return nil or a closed channel
	// Also this function may called several time and should return the same channel each time
	// and should not block
	Watch() <-chan map[string]interface{}
}

var o = &Onion{}

// Onion is a layer base configuration system
type Onion struct {
	lock sync.RWMutex

	delimiter string
	ll        []Layer

	// Loaded data
	data map[Layer]map[string]interface{}

	reload chan struct{}
}

func (o *Onion) watchLayer(ctx context.Context, l Layer) {
	c := l.Watch()
	if c == nil {
		return
	}
	for {
		select {
		case data, ok := <-c:
			if !ok {
				return
			}
			o.setLayerData(l, data, true)
		case <-ctx.Done():
			return
		}
	}
}

func (o *Onion) setLayerData(l Layer, data map[string]interface{}, watch bool) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.data == nil {
		o.data = make(map[Layer]map[string]interface{})
	}
	o.data[l] = data

	if !watch || o.reload == nil {
		return
	}

	close(o.reload)
	o.reload = nil
}

// AddLayersContext add a new layer to global config
func AddLayersContext(ctx context.Context, l ...Layer) {
	o.AddLayersContext(ctx, l...)
}

// AddLayersContext add new layers to the end of config layers. last layer is loaded after all other
// layer
func (o *Onion) AddLayersContext(ctx context.Context, l ...Layer) {
	if len(l) == 0 {
		return
	}
	o.lock.Lock()
	o.ll = append(o.ll, l...)
	o.lock.Unlock()

	for i := range l {
		o.setLayerData(l[i], l[i].Load(), false)
		go o.watchLayer(ctx, l[i])
	}
}

// AddLayers add a new layer to global config
func AddLayers(l ...Layer) {
	o.AddLayersContext(context.Background(), l...)
}

// AddLayers add new layers to onion
func (o *Onion) AddLayers(l ...Layer) {
	o.AddLayersContext(context.Background(), l...)
}

// GetDelimiter return the delimiter for nested key
func GetDelimiter() string {
	return o.GetDelimiter()
}

// GetDelimiter return the delimiter for nested key
func (o *Onion) GetDelimiter() string {
	if o.delimiter == "" {
		o.delimiter = "."
	}

	return o.delimiter
}

// SetDelimiter set the current delimiter on global config
func SetDelimiter(d string) {
	o.SetDelimiter(d)
}

// SetDelimiter set the current delimiter
func (o *Onion) SetDelimiter(d string) {
	o.delimiter = d
}

// ReloadWatch see onion.ReloadWatch
func ReloadWatch() <-chan struct{} {
	return o.ReloadWatch()
}

// ReloadWatch returns a channel to watch new layer data change, it just work for once, after the first change
// the channel will be changed to a new channel (the old channel will be closed to signal all listeners)
func (o *Onion) ReloadWatch() <-chan struct{} {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.reload == nil {
		o.reload = make(chan struct{})
	}

	return o.reload
}

// Get try to get the key from config layers
func Get(key string) (interface{}, bool) {
	return o.Get(key)
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
func GetIntDefault(key string, def int) int {
	return o.GetIntDefault(key, def)
}

// GetIntDefault return an int value from Onion, if the value is not exists or its not an
// integer , default is returned
func (o *Onion) GetIntDefault(key string, def int) int {
	return int(o.GetInt64Default(key, int64(def)))
}

// GetInt return an int value, if the value is not there, then it return zero value
func GetInt(key string) int {
	return o.GetInt(key)
}

// GetInt return an int value, if the value is not there, then it return zero value
func (o *Onion) GetInt(key string) int {
	return o.GetIntDefault(key, 0)
}

// GetInt64Default return an int64 value from Onion, if the value is not exists or if the value is not
// int64 then return the default
func GetInt64Default(key string, def int64) int64 {
	return o.GetInt64Default(key, def)
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
func GetInt64(key string) int64 {
	return o.GetInt64(key)
}

// GetInt64 return the int64 value from config, if its not there, return zero
func (o *Onion) GetInt64(key string) int64 {
	return o.GetInt64Default(key, 0)
}

// GetFloat32Default return an float32 value from Onion, if the value is not exists or its not a
// float32, default is returned
func GetFloat32Default(key string, def float32) float32 {
	return o.GetFloat32Default(key, def)
}

// GetFloat32Default return an float32 value from Onion, if the value is not exists or its not a
// float32, default is returned
func (o *Onion) GetFloat32Default(key string, def float32) float32 {
	return float32(o.GetFloat64Default(key, float64(def)))
}

// GetFloat32 return an float32 value, if the value is not there, then it returns zero value
func GetFloat32(key string) float32 {
	return o.GetFloat32(key)
}

// GetFloat32 return an float32 value, if the value is not there, then it returns zero value
func (o *Onion) GetFloat32(key string) float32 {
	return o.GetFloat32Default(key, 0)
}

// GetFloat64Default return an float64 value from Onion, if the value is not exists or if the value is not
// float64 then return the default
func GetFloat64Default(key string, def float64) float64 {
	return o.GetFloat64Default(key, def)
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
func GetFloat64(key string) float64 {
	return o.GetFloat64(key)
}

// GetFloat64 return the float64 value from config, if its not there, return zero
func (o *Onion) GetFloat64(key string) float64 {
	return o.GetFloat64Default(key, 0)
}

// GetStringDefault get a string from Onion. if the value is not exists or if tha value is not
// string, return the default
func GetStringDefault(key string, def string) string {
	return o.GetStringDefault(key, def)
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
func GetString(key string) string {
	return o.GetString(key)
}

// GetString is for getting an string from conig. if the key is not
func (o *Onion) GetString(key string) string {
	return o.GetStringDefault(key, "")
}

// GetBoolDefault return bool value from Onion. if the value is not exists or if tha value is not
// boolean, return the default
func GetBoolDefault(key string, def bool) bool {
	return o.GetBoolDefault(key, def)
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
func GetBool(key string) bool {
	return o.GetBool(key)
}

// GetBool is used to get a boolean value fro config, with false as default
func (o *Onion) GetBool(key string) bool {
	return o.GetBoolDefault(key, false)
}

// GetDurationDefault is a function to get duration from config. it support both
// string duration (like 1h3m2s) and integer duration
func GetDurationDefault(key string, def time.Duration) time.Duration {
	return o.GetDurationDefault(key, def)
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
func GetDuration(key string) time.Duration {
	return o.GetDuration(key)
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
func GetStringSlice(key string) []string {
	return o.GetStringSlice(key)
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

// LayersData is used to get all layers data at once
func (o *Onion) LayersData() []map[string]interface{} {
	o.lock.RLock()
	defer o.lock.RUnlock()

	res := make([]map[string]interface{}, 0, len(o.ll))
	for i := range o.ll {
		l := o.ll[i]
		res = append(res, o.data[l])
	}

	return res
}

// NewContext return a new Onion, context is used for watch
func NewContext(ctx context.Context, layers ...Layer) *Onion {
	o := &Onion{}

	o.AddLayersContext(ctx, layers...)
	return o
}

// New returns a new onion
func New(layers ...Layer) *Onion {
	return NewContext(context.Background(), layers...)
}
