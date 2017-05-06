// Package yamlloader is used to handle yaml file in Onion file/folder layer.
// for using this package, just import it
//
// 		import (
//			_ "github.com/fzerorubigd/onion/yamlloader"
//		)
//
// There is no need to do anything else, if you load a file with yaml/yml
// extension, the yaml loader is doing his job.
package yamlloader

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"gopkg.in/fzerorubigd/onion.v3"
)

type yamlLoader struct {
}

func (yl yamlLoader) SupportedEXT() []string {
	return []string{".yaml", ".yml"}
}

func (yl yamlLoader) Convert(r io.Reader) (map[string]interface{}, error) {
	data, _ := ioutil.ReadAll(r)
	ret := make(map[string]interface{})
	err := yaml.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func init() {
	onion.RegisterLoader(&yamlLoader{})
}
