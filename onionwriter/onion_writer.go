package onionwriter

import (
	"encoding/json"
	"io"

	"github.com/goraz/onion"
	"github.com/goraz/onion/helper"
)

// SerializeOnion try to serialize the onion into a json stream.
// TODO : Add more option, maybe support for more format. (Do we need more option on writing?)
func SerializeOnion(o *onion.Onion, w io.Writer) error {
	data := o.LayersData()

	mergedData := helper.MergeLayersData(data)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(mergedData)
}
