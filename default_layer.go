package onion

import (
	"errors"
	"strings"
)

// DefaultLayer is a layer to handle default value for layer.
type DefaultLayer interface {
	Layer
	// SetDefault set a default value for a key
	SetDefault(string, interface{}) error
	// GetDelimiter is used to get current delimiter for this layer. since
	// this layer needs to work with keys, the delimiter is needed
	GetDelimiter() string
	// SetDelimiter is used to set delimiter on this layer
	SetDelimiter(d string)
}

type defaultLayer struct {
	delimiter string
	data      map[string]interface{}
}

// Again, the case of two identical function and not convert one to another
func stringSetDefault(k []string, v interface{}, scope map[string]interface{}) error {
	if len(k) == 1 {
		// this is the key. just set it
		scope[k[0]] = v
		return nil
	}
	t, ok := scope[k[0]]
	if ok {
		// the key is already there. check if its another map?
		switch m := t.(type) {
		case map[string]interface{}:
			return stringSetDefault(k[1:], v, m)
		case map[interface{}]interface{}:
			return interfaceSetDefault(k[1:], v, m)
		default:
			return errors.New("the key is not a map")
		}
	}

	scope[k[0]] = make(map[string]interface{})
	return stringSetDefault(k[1:], v, scope[k[0]].(map[string]interface{}))
}

func interfaceSetDefault(k []string, v interface{}, scope map[interface{}]interface{}) error {
	if len(k) == 1 {
		// this is the key. just set it
		scope[k[0]] = v
		return nil
	}
	t, ok := scope[k[0]]
	if ok {
		// the key is already there. check if its another map?
		switch m := t.(type) {
		case map[string]interface{}:
			return stringSetDefault(k[1:], v, m)
		case map[interface{}]interface{}:
			return interfaceSetDefault(k[1:], v, m)
		default:
			return errors.New("the key is not a map")
		}
	}

	scope[k[0]] = make(map[string]interface{})
	return stringSetDefault(k[1:], v, scope[k[0]].(map[string]interface{}))
}

func (dl *defaultLayer) Load() (map[string]interface{}, error) {
	return dl.data, nil
}

func (dl *defaultLayer) SetDefault(k string, v interface{}) error {
	ka := strings.Split(k, dl.GetDelimiter())
	return stringSetDefault(ka, v, dl.data)
}

// GetDelimiter return the delimiter for nested key
func (dl defaultLayer) GetDelimiter() string {
	if dl.delimiter == "" {
		dl.delimiter = "."
	}

	return dl.delimiter
}

// SetDelimiter set the current delimiter
func (dl *defaultLayer) SetDelimiter(d string) {
	dl.delimiter = d
}

// NewDefaultLayer is used to return a default layer. should load this layer
// before any other layer, and before adding it, must add default value before
// adding this layer to onion.
func NewDefaultLayer() DefaultLayer {
	return &defaultLayer{
		delimiter: DefaultDelimiter,
		data:      make(map[string]interface{}),
	}
}
