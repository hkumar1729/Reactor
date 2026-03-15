package watcher

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func StartWatching(events chan struct{}) {
	w, err := fsnotify.NewWatcher()
	lastModified := make(map[string]time.Time)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	// add all directories
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			if _, ok := ignoreDirs[info.Name()]; ok {
				return filepath.SkipDir
			}
			return w.Add(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			// ignore unwanted operations
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}
			// handle new directories
			if event.Op&fsnotify.Create != 0 {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					w.Add(event.Name)
					continue
				}
			}
			info, err := os.Stat(event.Name)
			if err != nil || info.IsDir() {
				continue
			}
			lastTime, exists := lastModified[event.Name]
			if exists && info.ModTime().Equal(lastTime) {
				continue
			}
			lastModified[event.Name] = info.ModTime()
			log.Println("file changed:", event.Name)
			select {
			case events <- struct{}{}:
			default:
			}
		case err := <-w.Errors:
			log.Println("Error:", err)
		}
	}
}
