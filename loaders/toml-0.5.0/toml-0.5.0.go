// Package toml050loader is used to handle toml file in Onion file layer.
// for using this package, just import it
//
// 		import (
// 			_ "github.com/goraz/onion/loaders/toml"
// 		)
//
// There is no need to do anything else, if you load a file with toml
// extension, the toml loader is doing his job.
package toml050loader

import (
	"context"
	"io"

	"github.com/goraz/onion"
	"github.com/pelletier/go-toml"
)

type toml0_5_0Loader struct {
}

func (tl *toml0_5_0Loader) Decode(_ context.Context, r io.Reader) (map[string]interface{}, error) {
	config, err := toml.LoadReader(r)
	if err != nil {
		return nil, err
	}

	return config.ToMap(), nil
}

func init() {
	onion.RegisterDecoder(&toml0_5_0Loader{}, "toml")
}
