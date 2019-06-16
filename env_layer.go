package onion

import (
	"os"
	"strings"
)

// NewEnvLayer create new layer using the whitelist of environment values.
func NewEnvLayer(separator string, whiteList ...string) Layer {
	data := make(map[string]interface{})
	for i := range whiteList {
		if s := os.Getenv(whiteList[i]); s != "" {
			data[strings.ToLower(whiteList[i])] = s
		}
	}

	return NewMapLayerSeparator(data, separator)
}
