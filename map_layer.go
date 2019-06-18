package onion

// mapLayer is a layer based on maps. this layer can be used in other type of layers
type mapLayer struct {
	c chan map[string]interface{}
}

func (m *mapLayer) Load() <-chan map[string]interface{} {
	return m.c
}

// NewMapLayerSeparator returns a basic map layer, this layer is simply holds a map of values
func NewMapLayer(data map[string]interface{}) Layer {
	ret := &mapLayer{
		c: make(chan map[string]interface{}, 1),
	}
	ret.c <- data
	close(ret.c)

	return ret
}
