//go:build darwin || windows

package ui

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type KeyWatcher struct {
	watcher  *fsnotify.Watcher
	dirPath  string
	onChange func()
	stopCh   chan struct{}
	doneCh   chan struct{}
	mu       sync.Mutex
	stopped  bool
}

func NewKeyWatcher(dirPath string, onChange func()) (*KeyWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		absPath = dirPath
	}

	if err := os.MkdirAll(absPath, 0700); err != nil {
		_ = watcher.Close()
		return nil, err
	}

	if err := watcher.Add(absPath); err != nil {
		_ = watcher.Close()
		return nil, err
	}

	_ = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && path != absPath {
			_ = watcher.Add(path)
		}
		return nil
	})

	kw := &KeyWatcher{
		watcher:  watcher,
		dirPath:  absPath,
		onChange: onChange,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}

	go kw.listen()

	return kw, nil
}

func (kw *KeyWatcher) listen() {
	defer close(kw.doneCh)

	var debounceTimer *time.Timer
	var debounceCh <-chan time.Time

	for {
		select {
		case <-kw.stopCh:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return

		case <-debounceCh:
			debounceTimer = nil
			debounceCh = nil
			if kw.onChange != nil {
				kw.onChange()
			}

		case event, ok := <-kw.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create != 0 {
				if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
					_ = kw.watcher.Add(event.Name)
				}
			}

			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename|fsnotify.Chmod) != 0 {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.NewTimer(50 * time.Millisecond)
				debounceCh = debounceTimer.C
			}

		case err, ok := <-kw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("key watcher error: %v", err)
		}
	}
}

func (kw *KeyWatcher) Stop() {
	kw.mu.Lock()
	if kw.stopped {
		kw.mu.Unlock()
		return
	}
	kw.stopped = true
	kw.mu.Unlock()

	close(kw.stopCh)
	_ = kw.watcher.Close()
	<-kw.doneCh
}
