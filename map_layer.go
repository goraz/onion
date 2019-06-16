package onion

import (
	"sync"
)

// mapLayer is a layer based on maps. this layer can be used in other type of layers
type mapLayer struct {
	lock sync.RWMutex
	m    map[string]interface{}
}

func (ml *mapLayer) Load() (map[string]interface{}, error) {
	return ml.m, nil
}

// NewMapLayerSeparator returns a basic map layer, this layer is simply holds a map of values
func NewMapLayer(data map[string]interface{}) Layer {
	return &mapLayer{
		m:    data,
		lock: sync.RWMutex{},
	}
}
