package onion

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

var (
	decLock  sync.RWMutex
	decoders = map[string]Decoder{
		"json": &jsonDecoder{},
	}
)

// Decoder is a stream decoder to convert a stream into a map of config keys, json is supported out of
// the box
type Decoder interface {
	Decode(io.Reader) (map[string]interface{}, error)
}

type jsonDecoder struct {
}

func (jd *jsonDecoder) Decode(r io.Reader) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RegisterDecoder add a new decoder to the system, json is registered out of the box
func RegisterDecoder(typ string, dec Decoder) {
	decLock.Lock()
	defer decLock.Unlock()

	decoders[strings.ToLower(typ)] = dec
}

// NewStreamLayer returns a layer based on a io.Reader stream
func NewStreamLayer(stream io.Reader, format string) (Layer, error) {
	decLock.RLock()
	defer decLock.RUnlock()
	
	dec, ok := decoders[strings.ToLower(format)]
	if !ok {
		return nil, fmt.Errorf("there is no decoder for %q", format)
	}

	data, err := dec.Decode(stream)
	if err != nil {
		return nil, err
	}

	return NewMapLayer(data), nil
}
