// Package extraenv is the loader from the os env preceded by a prefix if set.
// for example for getting the test.path.key it check for {PREFIX_}TEST_PATH_KEY from the env
package extraenv

import (
	"os"
	"math"
	"strings"

	"gopkg.in/fzerorubigd/onion.v3"
)

type envLoader struct {
	prefix string
}

func (el *envLoader) Get(path ...string) (interface{}, bool) {
	p := el.prefix + strings.Repeat("_", int(math.Min(float64(len(el.prefix)), 1))) + strings.ToUpper(strings.Join(path, "_"))
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
