package onion

import (
	"fmt"
	"path/filepath"
)

type folderLayer struct {
	folder     string
	configName string

	loaded bool
	file   Layer
}

func (fl *folderLayer) IsLazy() bool {
	return false
}

func (fl *folderLayer) Load(d string, p ...string) (map[string]interface{}, error) {
	if fl.loaded {
		return fl.file.Load(d, p...)
	}

	files, err := filepath.Glob(fl.folder + "/" + fl.configName + ".*")
	if err != nil {
		return nil, err
	}
	for i := range files {
		// Try to load each file, until the first one is accepted
		fl.file = NewFileLayer(files[i])
		data, err := fl.file.Load(d, p...)
		if err == nil {
			fl.loaded = true
			return data, nil
		}
	}

	return nil, fmt.Errorf("there is no supported file in %s", fl.folder)
}

// NewFolderLayer return a new folder layer, this layer search in a folder for
// all supported file, and when it hit the first loadable file then simply return it
// the config name must not contain file extension
func NewFolderLayer(folder, configName string) Layer {
	// TODO : os specific seperator
	if folder[len(folder)-1:] != "/" {
		folder += "/"
	}

	return &folderLayer{folder, configName, false, nil}
}
