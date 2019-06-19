// Package tomlloader is used to handle toml file in Onion file/folder layer.
// for using this package, just import it
//
// 		import (
// 			_ "github.com/fzerorubigd/onion/tomlloader"
// 		)
//
// There is no need to do anything else, if you load a file with toml
// extension, the toml loader is doing his job.
package tomlloader

import (
	"context"
	"io"

	"github.com/BurntSushi/toml"

	"github.com/fzerorubigd/onion"
)

type tomlLoader struct {
}

func (tl *tomlLoader) Decode(_ context.Context, r io.Reader) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	_, err := toml.DecodeReader(r, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func init() {
	onion.RegisterDecoder(&tomlLoader{}, "toml")
}
