package main

import (
	"log"

	"github.com/howeyc/fsnotify"
)

func watcher(ch chan string, paths []string) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Setting up watcher: %v", err)
	}

	for _, path := range paths {
		err = w.Watch(path)
		if err != nil {
			log.Fatalf("Setting up watch path %v: %v", path, err)
		}
	}

	for {
		select {
		case e := <-w.Event:
			if !e.IsDelete() {
				ch <- e.Name
			}

		case err := <-w.Error:
			log.Fatalf("Watch error: %v", err)
		}
	}
}
