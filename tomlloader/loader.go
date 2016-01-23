// Package tomlloader is used to handle toml file in Onion file/folder layer.
// for using this package, just import it
//
// 		import (
//			_ "github.com/fzerorubigd/onion/tomlloader"
//		)
//
// There is no need to do anything else, if you load a file with toml
// extension, the toml loader is doing his job.
package tomlloader

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"gopkg.in/fzerorubigd/onion.v2"
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
