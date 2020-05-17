package directorylayer

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/goraz/onion"
	"github.com/skarademir/naturalsort"
)

// NewDirectoryLayer return a new directory layer.
// This layer search in a directory for all files with filesExtension extension
// and will use each of them as a file layer
func NewDirectoryLayer(directory, filesExtension string) (onion.Layer, error) {
	if directory[len(directory)-1:] != string(os.PathSeparator) {
		directory += string(os.PathSeparator)
	}

	fileNames := getFilesInOrder(directory, filesExtension)

	if len(fileNames) == 0 {
		return onion.NewMapLayer(nil), nil
	}

	layersData := make([]map[string]interface{}, 0)

	for _, fileName := range fileNames {
		layer, err := onion.NewFileLayer(fileName, nil)

		if err != nil {
			return nil, err
		}

		layersData = append(layersData, layer.Load())
	}

	return onion.NewMapLayer(layersData...), nil
}

func getFilesInOrder(directory, filesExtension string) []string {
	filePaths, err := filepath.Glob(directory + string(os.PathSeparator) + "*." + filesExtension)
	if err != nil {
		return nil
	}

	sort.Sort(naturalsort.NaturalSort(filePaths))

	return filePaths
}
