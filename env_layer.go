package onion

import (
	"os"
	"strings"
)

type envLoader struct {
	whiteList []string

	loaded bool
	data   map[string]interface{}
}

func (el *envLoader) Load() (map[string]interface{}, error) {
	if el.loaded {
		return el.data, nil
	}
	for _, env := range el.whiteList {
		v := os.Getenv(strings.ToUpper(env))
		if v != "" {
			el.data[env] = v
		}
	}
	el.loaded = true
	return el.data, nil
}

// NewEnvLayer create a environment loader. this loader accept a whitelist of allowed variables
// TODO : find a way to map env variable with different name
func NewEnvLayer(whiteList ...string) Layer {
	return &envLoader{whiteList, false, make(map[string]interface{})}
}
