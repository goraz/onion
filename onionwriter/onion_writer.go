package onionwriter

import (
	"encoding/json"
	"io"

	"github.com/goraz/onion"
	"github.com/mitchellh/mapstructure"
)

// SerializeOnion try to serialize the onion into a json stream.
func SerializeOnion(o *onion.Onion, w io.Writer) error {
	data := o.LayersData()

	mergedData := onion.NewMapLayer(data...)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(mergedData.Load())
}

// MergeLayersOnion is used to get all layers data merged into one
// Latest added overwrite previous ones.
func MergeLayersOnion(o *onion.Onion) map[string]interface{} {
	layersData := o.LayersData()

	return onion.NewMapLayer(layersData...).Load()
}

// DecodeOnion try to convert merged layers in the output structure.
// output must be a pointer to a map or struct.
func DecodeOnion(o *onion.Onion, output interface{}) error {
	merged := MergeLayersOnion(o)

	return mapstructure.Decode(merged, &output)
}
