package onion

import (
	"os"
	"strings"
)

func buildMap(m map[string]interface{}, v interface{}, k ...string) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}

	switch len(k) {
	case 0:
		return m
	case 1:
		m[k[0]] = v
		return m
	}
	d, _ := m[k[0]].(map[string]interface{})
	m[k[0]] = buildMap(d, v, k[1:]...)
	return m
}

// NewEnvLayer create new layer using the whitelist of environment values.
func NewEnvLayer(separator string, whiteList ...string) Layer {
	var data map[string]interface{}
	for i := range whiteList {
		if v, ok := os.LookupEnv(whiteList[i]); ok {
			keys := strings.Split(strings.ToLower(whiteList[i]), separator)
			data = buildMap(data, v, keys ...)
		}
	}

	return NewMapLayer(data)
}

// NewEnvLayerPrefix create new env layer, with all values with the same prefix
func NewEnvLayerPrefix(separator string, prefix string) Layer {
	var data map[string]interface{}
	pf := strings.ToUpper(prefix) + separator
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, pf) {
			k := strings.Trim(strings.Split(env, "=")[0], "\t\n ")
			ck := strings.ToLower(strings.TrimPrefix(k, pf))
			data = buildMap(data, os.Getenv(k), strings.Split(ck, separator) ...)
		}
	}

	return NewMapLayer(data)
}
