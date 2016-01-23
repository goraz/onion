package extraenv

import (
	"os"
	"strings"

	"gopkg.in/fzerorubigd/onion.v2"
)

type envLoader struct {
	prefix string
}

func (el *envLoader) IsLazy() bool {
	return true
}

func createNestedMap(v interface{}, p ...string) interface{} {
	if len(p) == 0 {
		return v
	}
	return map[string]interface{}{p[0]: createNestedMap(v, p[1:]...)}
}

func (el *envLoader) Load(d string, path ...string) (map[string]interface{}, error) {
	if len(path) == 0 {
		return nil, nil
	}

	p := el.prefix + "_" + strings.ToUpper(strings.Join(path, "_"))
	v := os.Getenv(p)
	m := make(map[string]interface{})
	if v != "" && len(p) > 0 {
		m = createNestedMap(v, path...).(map[string]interface{})
	}
	return m, nil
}

// NewExtraEnvLayer create a environment loader. this layer is base on the influxdb config
func NewExtraEnvLayer(prefix string) onion.Layer {
	return &envLoader{
		prefix: strings.ToUpper(prefix),
	}
}
