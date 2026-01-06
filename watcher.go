package main

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher
	done      chan struct{}
}

func NewWatcher(filePath string, callback func(string)) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = fsWatcher.Add(filePath)
	if err != nil {
		_ = fsWatcher.Close()
		return nil, err
	}

	w := &Watcher{
		fsWatcher: fsWatcher,
		done:      make(chan struct{}),
	}

	go func() {
		for {
			select {
			case event, ok := <-fsWatcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					content, err := os.ReadFile(filePath)
					if err == nil {
						callback(string(content))
					}
				}
			case <-w.done:
				return
			}
		}
	}()

	return w, nil
}

func (w *Watcher) Close() error {
	close(w.done)
	return w.fsWatcher.Close()
}
