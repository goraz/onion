package onion

import (
	"strings"
	"sync"
)

// mapLayer is a layer based on maps. this layer can be used in other type of layers
type mapLayer struct {
	lock sync.RWMutex
	m    map[string]interface{}
	sep  string
}

func searchStringMap(m map[string]interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(m, path[1:]...)
	case map[interface{}]interface{}:
		return searchInterfaceMap(m, path[1:]...)
	}
	return nil, false
}

func searchInterfaceMap(m map[interface{}]interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(m, path[1:]...)
	case map[interface{}]interface{}:
		return searchInterfaceMap(m, path[1:]...)
	}
	return nil, false
}

// Reload is called when the onion is going to reload itself. the MapLayer is not reloadable,
// its internal map is always the same
func (ml *mapLayer) Reload() (bool, error) {
	return false, nil
}

func (ml *mapLayer) Get(keys ...string) (interface{}, bool) {
	ml.lock.RLock()
	defer ml.lock.RUnlock()

	if ml.sep != "" {
		key := strings.Join(keys, ml.sep)
		if v, ok := ml.m[key]; ok {
			return v, true
		}
	}
	return searchStringMap(ml.m, keys...)
}

// NewMapLayerSeparator returns a basic map layer, this layer is simply holds a map of values and
// search inside the map recursively.
// Also if separator is not empty string, the first level keys with the exact math returned instead
// of recursive search
func NewMapLayerSeparator(data map[string]interface{}, separator string) Layer {
	return &mapLayer{
		m:    data,
		lock: sync.RWMutex{},
		sep:  separator,
	}
}

// NewMapLayer create a new map layer
func NewMapLayer(data map[string]interface{}) Layer {
	return NewMapLayerSeparator(data, "")
}
