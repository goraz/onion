package onion

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/goraz/onion/utils"
	"github.com/skarademir/naturalsort"
)

// NewFolderLayer return a new folder layer.
// This layer search in a folder for all files with filesExtension extension
// and will use each of them as a file layer
func NewFolderLayer(folder, filesExtension string) (Layer, error) {
	if folder[len(folder)-1:] != string(os.PathSeparator) {
		folder += string(os.PathSeparator)
	}

	fileNames := getFilesInOrder(folder, filesExtension)
	layersData := make([]map[string]interface{}, 0)

	for _, fileName := range fileNames {
		layer, err := NewFileLayer(fileName, nil)

		if err != nil {
			return nil, err
		}

		layersData = append(layersData, layer.Load())
	}

	return NewMapLayer(utils.MergeLayersData(layersData)), nil
}

func getFilesInOrder(folder, filesExtension string) []string {
	filePaths, err := filepath.Glob(folder + string(os.PathSeparator) + "*." + filesExtension)
	if err != nil {
		return nil
	}

	sort.Sort(naturalsort.NaturalSort(filePaths))

	return filePaths
}
