//go:build darwin

package ui

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type ConfigWatcher struct {
	watcher  *fsnotify.Watcher
	filePath string
	onChange func()
	stopCh   chan struct{}
	doneCh   chan struct{}
	mu       sync.Mutex
	stopped  bool
}

func NewConfigWatcher(filePath string, onChange func()) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}

	dir := filepath.Dir(absPath)
	if err := watcher.Add(dir); err != nil {
		_ = watcher.Close()
		return nil, err
	}

	// Also watch the file directly if it exists
	_ = watcher.Add(absPath)

	cw := &ConfigWatcher{
		watcher:  watcher,
		filePath: absPath,
		onChange: onChange,
		stopCh:   make(chan struct{}),
		doneCh:   make(chan struct{}),
	}

	go cw.listen()

	return cw, nil
}

func (cw *ConfigWatcher) listen() {
	defer close(cw.doneCh)

	cleanTarget := filepath.Clean(cw.filePath)
	targetBase := filepath.Base(cw.filePath)

	var debounceTimer *time.Timer
	var debounceCh <-chan time.Time

	for {
		select {
		case <-cw.stopCh:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return

		case <-debounceCh:
			debounceTimer = nil
			debounceCh = nil
			if cw.onChange != nil {
				cw.onChange()
			}

		case event, ok := <-cw.watcher.Events:
			if !ok {
				return
			}

			eventClean := filepath.Clean(event.Name)
			eventBase := filepath.Base(event.Name)

			if eventClean == cleanTarget || eventBase == targetBase {
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Chmod) != 0 {
					// Re-add the target file in case of rename/recreation
					_ = cw.watcher.Add(cw.filePath)

					// Debounce multiple rapid events (e.g. Write + Chmod)
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.NewTimer(50 * time.Millisecond)
					debounceCh = debounceTimer.C
				}
			}

		case err, ok := <-cw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("config watcher error: %v", err)
		}
	}
}

func (cw *ConfigWatcher) Stop() {
	cw.mu.Lock()
	if cw.stopped {
		cw.mu.Unlock()
		return
	}
	cw.stopped = true
	cw.mu.Unlock()

	close(cw.stopCh)
	_ = cw.watcher.Close()
	<-cw.doneCh
}
