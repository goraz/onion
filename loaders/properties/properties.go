// Package properties is used to handle properties file in Onion file layer.
// for using this package, just import it
//
// 		import (
// 			_ "github.com/goraz/onion/loaders/properties"
// 		)
//
// There is no need to do anything else, if you load a file with toml
// extension, the toml loader is doing his job.
package properties

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/goraz/onion"
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
