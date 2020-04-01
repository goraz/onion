package onion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// Cipher is used to decrypt data on loading
type Cipher interface {
	Decrypt(io.Reader) ([]byte, error)
}

// Decoder is a stream decoder to convert a stream into a map of config keys, json is supported out of
// the box
type Decoder interface {
	Decode(context.Context, io.Reader) (map[string]interface{}, error)
}

type jsonDecoder struct {
}

func decrypt(c Cipher, r io.Reader) (io.Reader, error) {
	if c == nil {
		return r, nil
	}
	b, err := c.Decrypt(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func (jd *jsonDecoder) Decode(_ context.Context, r io.Reader) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// RegisterDecoder add a new decoder to the system, json is registered out of the box
func RegisterDecoder(dec Decoder, formats ...string) {
	decLock.Lock()
	defer decLock.Unlock()

	for _, format := range formats {
		format := strings.ToLower(format)

		_, alreadyExists := decoders[format]
		if alreadyExists {
			log.Fatalf("decoder for format %q is already registered: you can have only one", format)
		}

		decoders[format] = dec
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
	c      chan map[string]interface{}
	cipher Cipher
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
	dr, err := decrypt(sl.cipher, r)
	if err != nil {
		return err
	}

	data, err := dec.Decode(ctx, dr)
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

// NewStreamLayerContext try to create a layer based on a stream, the format should be a registered
// format (see RegisterDecoder) and if the Cipher is not nil, it pass data to cipher first.
// A nil cipher is accepted as plain cipher
func NewStreamLayerContext(ctx context.Context, r io.Reader, format string, c Cipher) (Layer, error) {
	if r == nil {
		return nil, fmt.Errorf("nil stream")
	}
	sl := &streamLayer{
		c:      make(chan map[string]interface{}),
		cipher: c,
	}

	err := sl.Reload(ctx, r, format)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

// NewStreamLayer create new stream layer, see the NewStreamLayerContext
func NewStreamLayer(r io.Reader, format string, c Cipher) (Layer, error) {
	return NewStreamLayerContext(context.Background(), r, format, c)
}

// NewFileLayerContext create a new file layer. it choose the format base on the extension
func NewFileLayerContext(ctx context.Context, path string, c Cipher) (Layer, error) {
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return NewStreamLayerContext(ctx, f, ext, c)
}

// NewFileLayer create a new file layer. it choose the format base on the extension
func NewFileLayer(path string, c Cipher) (Layer, error) {
	return NewFileLayerContext(context.Background(), path, c)
}
