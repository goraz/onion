package structlayer

import (
	"github.com/goraz/onion"
	"github.com/mitchellh/mapstructure"
)

// NewStructLayer returns a layer based on a structure.
func NewStructLayer(s interface{}) (onion.Layer, error) {
	mp := make(map[string]interface{})
	if err := mapstructure.Decode(s, &mp); err != nil {
		return nil, err
	}

	return onion.NewMapLayer(mp), nil
}
