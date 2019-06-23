package onion

// mapLayer is a layer based on maps. this layer can be used in other type of layers
type mapLayer struct {
	data map[string]interface{}
}

func (m *mapLayer) Load() map[string]interface{} {
	return m.data
}

func (m *mapLayer) Watch() <-chan map[string]interface{} {
	return nil
}

// NewMapLayer returns a basic map layer, this layer is simply holds a map of values
func NewMapLayer(data map[string]interface{}) Layer {
	ret := &mapLayer{
		data: data,
	}

	return ret
}
