package filewatchlayer

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/fzerorubigd/onion"
)

type streamReload interface {
	Reload(context.Context, io.Reader, string) error
}

func reload(ctx context.Context, path string, fl streamReload, ext string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return fl.Reload(ctx, f, ext)
}

// NewFileWatchLayerContext create a file layer with automatic fswatch.
// it reloads if the file content has changed, also the watch finish with the context
// a non-nil ciper is used to load encrypted file, nil means plain file
func NewFileWatchLayerContext(ctx context.Context, path string, c onion.Cipher) (onion.Layer, error) {
	l, err := onion.NewFileLayerContext(ctx, path, c)
	if err != nil {
		return nil, err
	}
	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	sl := l.(streamReload)
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := watch.Add(path); err != nil {
		return nil, err
	}

	go func() {
		defer func() { _ = watch.Close() }()
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watch.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := reload(ctx, path, sl, ext); err != nil {
					log.Println("error:", err) // Better log support
				}
			}
		case err, ok := <-watch.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}()

	return l, nil
}

// NewFileWatchLayer create a file layer with automatic fswatch.
// it reloads if the file content has changed
func NewFileWatchLayer(path string, c onion.Cipher) (onion.Layer, error) {
	return NewFileWatchLayerContext(context.Background(), path, c)
}
