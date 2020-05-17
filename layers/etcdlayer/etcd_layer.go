// Package etcdlayer is a layer to manage a configuration on a key inside the etcd, it watches the change on the key
package etcdlayer

import (
	"bytes"
	"context"
	"io"
	"log"
	"time"

	"github.com/goraz/onion"
	goetcd "go.etcd.io/etcd/client"
)

type streamReload interface {
	Reload(context.Context, io.Reader, string) error
}

func getWithContext(ctx context.Context, api goetcd.KeysAPI, key string) (io.Reader, error) {
	resp, err := api.Get(ctx, key, nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader([]byte(resp.Node.Value)), nil
}

func watchWithContext(ctx context.Context, api goetcd.KeysAPI, key string) <-chan []byte {
	respChan := make(chan []byte)
	go func() {
		watcher := api.Watcher(key, nil)
		for {
			resp, err := watcher.Next(ctx)
			if err != nil {
				log.Println("error:", err) // Better log support
				time.Sleep(time.Second * 5)
				continue
			}
			respChan <- []byte(resp.Node.Value)
		}
	}()
	return respChan
}

// NewEtcdLayerContext reads config from a etcd key, it should encode with one of the know formats and
// optionally can be encrypted using cipher.
func NewEtcdLayerContext(ctx context.Context, key string, format string, endPoints []string, c onion.Cipher) (onion.Layer, error) {
	cl, err := goetcd.New(goetcd.Config{
		Endpoints: endPoints,
	})
	if err != nil {
		return nil, err
	}
	api := goetcd.NewKeysAPI(cl)
	buf, err := getWithContext(ctx, api, key)
	if err != nil {
		return nil, err
	}

	l, err := onion.NewStreamLayerContext(ctx, buf, format, c)
	if err != nil {
		return nil, err
	}

	sl := l.(streamReload)

	go func() {
		watch := watchWithContext(ctx, api, key)
		for {
			select {
			case <-ctx.Done():
				return
			case b := <-watch:
				if err := sl.Reload(ctx, bytes.NewReader(b), format); err != nil {
					log.Println("error:", err) // Better log support
				}
			}
		}
	}()

	return l, nil
}

// NewEtcdLayer creates a new etcd layer
func NewEtcdLayer(key string, format string, endPoints []string, c onion.Cipher) (onion.Layer, error) {
	return NewEtcdLayerContext(context.Background(), key, format, endPoints, c)
}
