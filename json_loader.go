package onion

import (
	"encoding/json"
	"io"
)

type jsonLoader struct {
}

func (jl jsonLoader) SupportedEXT() []string {
	return []string{".json"}
}

func (jl jsonLoader) Convert(r io.Reader) (map[string]interface{}, error) {
	dec := json.NewDecoder(r)

	ret := make(map[string]interface{})
	err := dec.Decode(&ret)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func init() {
	RegisterLoader(&jsonLoader{})
}
