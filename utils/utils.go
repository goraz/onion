package utils

// MergeLayersData is an helper function to merge more layers in one.
// Following slice order, a previous layer key is overriden by an equal key in
// next layer.
func MergeLayersData(layers []map[string]interface{}) map[string]interface{} {
	mergedLayer := layers[len(layers)-1]
	layers = layers[:len(layers)-1]

	for i := len(layers) - 1; i >= 0; i-- {
		mergedLayer = mergeKeys(mergedLayer, layers[i])
	}

	return mergedLayer
}

// mergeKeys recursively merge right into left, never replacing any key that already exists in left
func mergeKeys(left, right map[string]interface{}) map[string]interface{} {
	if left == nil {
		return right
	}

	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			_, leftValIsAMap := leftVal.(map[string]interface{})
			_, rightValIsAMap := leftVal.(map[string]interface{})

			if leftValIsAMap && rightValIsAMap {
				left[key] = mergeKeys(leftVal.(map[string]interface{}), rightVal.(map[string]interface{}))
			}
		} else {
			left[key] = rightVal
		}
	}

	return left
}
