package directorywatchlayer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/goraz/onion"
	"github.com/skarademir/naturalsort"
)

var (
	ErrNotDir             = errors.New("path doesn't point to a directory")
	ErrReloadNotSupported = errors.New("layer doesn't support reload")
)

// NewDirectoryWatchLayerContext watch for changes of existing files TODO: Watch new files
func NewDirectoryWatchLayerContext(
	ctx context.Context,
	dir string,
	cipher onion.Cipher,
	extensions ...string,
) ([]onion.Layer, error) {
	dir = filepath.Clean(dir)

	if fs, err := os.Stat(dir); nil != err {
		return nil, err
	} else if !fs.IsDir() {
		return nil, ErrNotDir
	}

	files, errList := directoryListByExtensions(dir, extensions...)
	if nil != errList {
		return nil, errList
	}

	pathToLayerIndex := make(map[string]int)
	layers := make([]onion.Layer, len(files))
	for k, path := range files {
		if l, err := onion.NewFileLayerContext(ctx, path, cipher); nil == err {
			layers[k] = l
			pathToLayerIndex[path] = k
		} else {
			return nil, err
		}
	}

	watcher, errInit := fsnotify.NewWatcher()
	if nil != errInit {
		return nil, errInit
	}

	if err := watcher.Add(dir); nil != err {
		_ = watcher.Close()

		return nil, err
	}

	go func() {
		defer func() { _ = watcher.Close() }()

		for {
			select {
			case <-ctx.Done():
				return

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if fsnotify.Write == event.Op&fsnotify.Write {
					<-time.After(time.Second) // sometime it triggers before the complete write TODO: find a solution (not hack)

					_ = reloadLayer(ctx, layers[pathToLayerIndex[event.Name]], event.Name)
				}

			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	return layers, nil
}

func NewDirectoryWatchLayer(dir string, cipher onion.Cipher, extensions ...string) ([]onion.Layer, error) {
	return NewDirectoryWatchLayerContext(context.Background(), dir, cipher, extensions...)
}

func directoryListByExtensions(dir string, extensions ...string) ([]string, error) {
	patterns := make([]string, len(extensions))
	if 0 == len(patterns) {
		patterns = append(patterns, "*")
	} else {
		for k, ext := range extensions {
			patterns[k] = fmt.Sprintf("*.%s", ext)
		}
	}

	list := make([]string, 0)
	for _, pattern := range patterns {
		if paths, err := filepath.Glob(fmt.Sprintf("%s%c%s", dir, os.PathSeparator, pattern)); nil == err {
			list = append(list, paths...)
		} else {
			return nil, err
		}
	}

	sort.Sort(naturalsort.NaturalSort(list))

	return list, nil
}

type streamReload interface {
	Reload(context.Context, io.Reader, string) error
}

func reloadLayer(ctx context.Context, layer onion.Layer, path string) error {
	sl, ok := layer.(streamReload)
	if !ok {
		return ErrReloadNotSupported
	}

	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = fh.Close() }()

	ext := strings.TrimPrefix(filepath.Ext(path), ".")

	return sl.Reload(ctx, fh, ext)
}
