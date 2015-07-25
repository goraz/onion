package tomlloader

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"github.com/fzerorubigd/onion"
)

type tomlLoader struct {
}

func (tl tomlLoader) SupportedEXT() []string {
	return []string{".toml"}
}

func (tl tomlLoader) Convert(r io.Reader) (map[string]interface{}, error) {
	data, _ := ioutil.ReadAll(r)
	ret := make(map[string]interface{})
	err := toml.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func init() {
	onion.RegisterLoader(&tomlLoader{})
}
