// Package yamlloader is used to handle yaml file in Onion stream layer.
// for using this package, just import it
//
// 		import (
// 			_ "github.com/fzerorubigd/onion/yamlloader"
// 		)
//
// There is no need to do anything else, if you load a file with yaml/yml
// extension, the yaml loader is doing his job.
package yamlloader

import (
	"io"

	"gopkg.in/yaml.v2"

	"github.com/fzerorubigd/onion"
)

type yamlLoader struct {
}

func (yl yamlLoader) Decode(r io.Reader) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	err := yaml.NewDecoder(r).Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func init() {
	onion.RegisterDecoder(&yamlLoader{}, "yml", "yaml")
}
