package watcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func StartWatching(events chan struct{}) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

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
