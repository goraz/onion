package onion

func searchStringMap(m map[string]interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(m, path[1:]...)
	case map[interface{}]interface{}:
		return searchInterfaceMap(m, path[1:]...)
	}
	return nil, false
}

func searchInterfaceMap(m map[interface{}]interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	v, ok := m[path[0]]
	if !ok {
		return nil, false
	}

	if len(path) == 1 {
		return v, true
	}

	switch m := v.(type) {
	case map[string]interface{}:
		return searchStringMap(m, path[1:]...)
	case map[interface{}]interface{}:
		return searchInterfaceMap(m, path[1:]...)
	}
	return nil, false
}
