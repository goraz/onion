package onionwriter

import (
	"encoding/json"
	"io"

	"github.com/fzerorubigd/onion"
	"github.com/imdario/mergo"
)

// SerializeOnion try to serialize the onion into a json stream.
// TODO : Add more option, maybe support for more format. (Do we need more option on writing?)
func SerializeOnion(o *onion.Onion, w io.Writer) error {
	data := o.LayersData()
	res := make(map[string]interface{})
	for i := range data {
		err := mergo.Merge(&res, data[i], mergo.WithOverride)
		if err != nil {
			return err
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(res)
}
