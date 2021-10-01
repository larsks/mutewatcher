package main

import "sync"

type (
	Watcher struct {
		sync.Mutex
		watchedDevices []string
	}
)

func (w *Watcher) Watch() {
	for {
		w.Lock()
		w.Unlock()
	}
}

func (w *Watcher) Register(name string) {
	w.Lock()
	defer w.Unlock()
	for _, dev := range w.watchedDevices {
		if dev == name {
			return
		}
	}

	w.watchedDevices = append(w.watchedDevices, name)
}

func (w *Watcher) Unregister(name string) error {
	w.Lock()
	defer w.Unlock()
	return nil
}
