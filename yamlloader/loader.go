package yamlloader

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/fzerorubigd/onion"
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
