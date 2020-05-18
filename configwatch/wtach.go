package configwatch

import (
	"context"
	"sync"
	"time"

	"github.com/goraz/onion"
)

var (
	w = RefWatch{}
)

type variable struct {
	key string
	ref interface{}
	def interface{}
}

// RefWatch is a class to watch over the onion instance and change the registered variables
type RefWatch struct {
	refLock    sync.RWMutex
	references []variable
}

func (rw *RefWatch) addRef(key string, ref interface{}, def interface{}) {
	rw.refLock.Lock()
	defer rw.refLock.Unlock()

	rw.references = append(rw.references, variable{key: key, ref: ref, def: def})
}

// RegisterInt return an int variable and set the value when the config is loaded
func (rw *RefWatch) RegisterInt(key string, def int) Int {
	var v = int64(def)
	rw.addRef(key, &v, def)

	return intHolder{value: &v}
}

// RegisterInt return an int variable and set the value when the config is loaded
func RegisterInt(key string, def int) Int {
	return w.RegisterInt(key, def)
}

// RegisterInt64 return an int64 variable and set the value when the config is loaded
func (rw *RefWatch) RegisterInt64(key string, def int64) Int {
	var v = def
	rw.addRef(key, &v, def)

	return intHolder{value: &v}
}

// RegisterInt64 return an int64 variable and set the value when the config is loaded
func RegisterInt64(key string, def int64) Int {
	return w.RegisterInt64(key, def)
}

// RegisterString return an string variable and set the value when the config is loaded
func (rw *RefWatch) RegisterString(key string, def string) String {
	var v = def
	rw.addRef(key, &v, def)

	return stringHolder{value: &v}
}

// RegisterString return an string variable and set the value when the config is loaded
func RegisterString(key string, def string) String {
	return w.RegisterString(key, def)
}

// RegisterFloat64 return an float64 variable and set the value when the config is loaded
func (rw *RefWatch) RegisterFloat64(key string, def float64) Float {
	var v = def
	rw.addRef(key, &v, def)

	return floatHolder{value: &v}
}

// RegisterFloat64 return an float64 variable and set the value when the config is loaded
func RegisterFloat64(key string, def float64) Float {
	return w.RegisterFloat64(key, def)
}

// RegisterFloat32 return an float32 variable and set the value when the config is loaded
func (rw *RefWatch) RegisterFloat32(key string, def float32) Float {
	var v = float64(def)
	rw.addRef(key, &v, def)

	return floatHolder{value: &v}
}

// RegisterFloat32 return an float32 variable and set the value when the config is loaded
func RegisterFloat32(key string, def float32) Float {
	return w.RegisterFloat32(key, def)
}

// RegisterBool return an bool variable and set the value when the config is loaded
func (rw *RefWatch) RegisterBool(key string, def bool) Bool {
	var v = def
	rw.addRef(key, &v, def)

	return boolHolder{value: &v}
}

// RegisterBool return an bool variable and set the value when the config is loaded
func RegisterBool(key string, def bool) Bool {
	return w.RegisterBool(key, def)
}

// RegisterDuration return an duration variable and set the value when the config is loaded
func (rw *RefWatch) RegisterDuration(key string, def time.Duration) Int {
	var v = int64(def)
	rw.addRef(key, &v, def)

	return intHolder{value: &v}
}

// RegisterDuration return an duration variable and set the value when the config is loaded
func RegisterDuration(key string, def time.Duration) Int {
	return w.RegisterDuration(key, def)
}

func (rw *RefWatch) watchLoop(o *onion.Onion) {
	rw.refLock.RLock()

	// Make sure all variables are locked
	// TODO : lock per onion instance
	globalLock.Lock()

	for i := range rw.references {
		switch def := rw.references[i].def.(type) {
		case int:
			v := o.GetInt64Default(rw.references[i].key, int64(def))
			t := rw.references[i].ref.(*int64)
			*t = v
		case int64:
			v := o.GetInt64Default(rw.references[i].key, def)
			t := rw.references[i].ref.(*int64)
			*t = v
		case string:
			v := o.GetStringDefault(rw.references[i].key, def)
			t := rw.references[i].ref.(*string)
			*t = v
		case float32:
			v := o.GetFloat64Default(rw.references[i].key, float64(def))
			t := rw.references[i].ref.(*float64)
			*t = v
		case float64:
			v := o.GetFloat64Default(rw.references[i].key, def)
			t := rw.references[i].ref.(*float64)
			*t = v
		case bool:
			v := o.GetBoolDefault(rw.references[i].key, def)
			t := rw.references[i].ref.(*bool)
			*t = v
		case time.Duration:
			v := o.GetDurationDefault(rw.references[i].key, def)
			t := rw.references[i].ref.(*int64)
			*t = int64(v)
		}
	}
	rw.refLock.RUnlock()
	globalLock.Unlock()

}

// Watch get an onion and watch over it for changes in the layers.
func (rw *RefWatch) Watch(ctx context.Context, o *onion.Onion) <-chan struct{} {
	rw.watchLoop(o)
	w := make(chan struct{})
	go func() {
		for {
			ch := o.ReloadWatch()
			select {
			case <-ctx.Done():
				return
			case <-ch:
				rw.watchLoop(o)
				select {
				case w <- struct{}{}:
				default:
				}
			}
		}
	}()

	return w
}

// Watch get an onion and watch over it for changes in the layers.
func Watch(o *onion.Onion) <-chan struct{} {
	return WatchContext(context.Background(), o)
}

// WatchContext get an onion and watch over it for changes in the layers.
func WatchContext(ctx context.Context, o *onion.Onion) <-chan struct{} {
	return w.Watch(ctx, o)
}
