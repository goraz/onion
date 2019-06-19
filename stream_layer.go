package onion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	Decode(context.Context, io.Reader) (map[string]interface{}, error)
}

type jsonDecoder struct {
}

func (jd *jsonDecoder) Decode(_ context.Context, r io.Reader) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RegisterDecoder add a new decoder to the system, json is registered out of the box
func RegisterDecoder(dec Decoder, typ ...string) {
	decLock.Lock()
	defer decLock.Unlock()

	for i := range typ {
		decoders[strings.ToLower(typ[i])] = dec
	}
}

// GetDecoder returns the decoder based on its name, it may returns nil if the decoder is not
// registered
func GetDecoder(format string) Decoder {
	decLock.RLock()
	defer decLock.RUnlock()

	return decoders[strings.ToLower(format)]
}

type streamLayer struct {
	c chan map[string]interface{}
}

func (sl *streamLayer) Load() map[string]interface{} {
	return <-sl.c
}

func (sl *streamLayer) Watch() <-chan map[string]interface{} {
	return sl.c
}

func (sl *streamLayer) Reload(ctx context.Context, r io.Reader, format string) error {
	dec := GetDecoder(format)
	if dec == nil {
		return fmt.Errorf("format %q is not registered", format)
	}
	data, err := dec.Decode(ctx, r)
	if err != nil {
		return err
	}

	go func() {
		select {
		case sl.c <- data:
		case <-ctx.Done():
		}
	}()

	return nil
}

func NewStreamLayerContext(ctx context.Context, r io.Reader, format string) (Layer, error) {
	sl := &streamLayer{
		c: make(chan map[string]interface{}),
	}

	err := sl.Reload(ctx, r, format)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func NewStreamLayer(r io.Reader, format string) (Layer, error) {
	return NewStreamLayerContext(context.Background(), r, format)
}

func NewFileLayerContext(ctx context.Context, path string) (Layer, error) {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return NewStreamLayerContext(ctx, f, ext)
}

func NewFileLayer(path string) (Layer, error) {
	return NewFileLayerContext(context.Background(), path)
}
