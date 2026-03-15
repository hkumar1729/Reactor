package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func StartWatching(events chan struct{}) {
	w, err := fsnotify.NewWatcher()
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
			log.Println("watching", path)
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
			if event.Has(fsnotify.Write) {
				log.Println("file modified:", event.Name)
				events <- struct{}{}
			}
		case err := <-w.Errors:
			log.Println("Error:", err)
		}
	}
}
