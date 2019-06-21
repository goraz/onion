package properties

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/fzerorubigd/onion"
	"github.com/magiconair/properties"
)

type propertiesLoader struct {
}

func (tl *propertiesLoader) Decode(_ context.Context, r io.Reader) (map[string]interface{}, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	p, err := properties.Load(data, properties.UTF8)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	for k, v := range p.Map() {
		ret[k] = v
	}

	return ret, nil
}

func init() {
	onion.RegisterDecoder(&propertiesLoader{}, "properties", "props")
}
