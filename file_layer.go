package onion

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var loaders map[string]FileLoader

// FileLoader is an interface to handle load config from a file
type FileLoader interface {
	// SupportedEXT Must return the list of supported ext for this loader interface
	SupportedEXT() []string
	// Convert is for translating the file data into config structure.
	Convert(io.Reader) (map[string]interface{}, error)
}

type fileLayer struct {
	file   string
	loaded bool
	data   map[string]interface{}
}

// Load a file. also save it's data so the next time it can simply return it
// may be I should remove cache?
func (fl *fileLayer) Load() (map[string]interface{}, error) {
	if fl.loaded {
		return fl.data, nil
	}
	ext := strings.ToLower(filepath.Ext(fl.file))
	l, ok := loaders[ext]
	if !ok {
		return nil, fmt.Errorf("no registered loader for ext %s", ext)
	}

	f, err := os.Open(fl.file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	fl.data, err = l.Convert(f)
	fl.loaded = err == nil

	return fl.data, err
}

// RegisterLoader must be called to register a type loader, this function is only available with
// file and folder loaders.
func RegisterLoader(l FileLoader) {
	for _, ext := range l.SupportedEXT() {
		loaders[strings.ToLower(ext)] = l
	}
}

// NewFileLayer initialize a new file layer. its for a single file, and the file ext
// is the key for loader to load a correct loader. if you want to scan a directory,
// use the folder loader.
func NewFileLayer(file string) Layer {
	return &fileLayer{
		file:   file,
		loaded: false,
		data:   nil,
	}
}

func init() {
	loaders = make(map[string]FileLoader)
}
