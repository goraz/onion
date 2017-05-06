// Package consulloader is a simple loader to handle Reading from consul.
// its a simple read only functionality.
package consulloader

import (
	"path/filepath"

	"github.com/hashicorp/consul/api"
	"gopkg.in/fzerorubigd/onion.v3"
)

type layer struct {
	client *api.KV
	prefix string
}

func (l layer) Get(path ...string) (interface{}, bool) {
	p := filepath.Join(l.prefix, filepath.Join(path...))
	kv, _, err := l.client.Get(p, nil)
	if err != nil {
		// I don't think this is correct. may be panic?
		return nil, false
	}
	if kv == nil {
		return nil, false
	}
	return string(kv.Value), true
}

// NewConsulLayer create a new lazy layer from a consul client
func NewConsulLayer(c *api.Client, prefix string) onion.LazyLayer {
	return &layer{
		client: c.KV(),
		prefix: prefix,
	}
}
