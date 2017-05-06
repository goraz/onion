// Package extraenv is the loader from the os env using prefix.
// for example for getting the test.path.key it check for PREFIX_TEST_PATH_KEY from the env
package extraenv

import (
	"os"
	"strings"

	"gopkg.in/fzerorubigd/onion.v3"
)

type envLoader struct {
	prefix string
}

func (el *envLoader) Get(path ...string) (interface{}, bool) {
	p := el.prefix + "_" + strings.ToUpper(strings.Join(path, "_"))
	v := os.Getenv(p)
	if v != "" && len(p) > 0 {
		return v, true
	}
	return nil, false
}

// NewExtraEnvLayer create a environment loader. this layer is base on the influxdb config
func NewExtraEnvLayer(prefix string) onion.LazyLayer {
	return &envLoader{
		prefix: strings.ToUpper(prefix),
	}
}
