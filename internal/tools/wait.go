package tools

import "sync"

// WaitGroupWrapper wrap sync.WaitGroup
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap wrap
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
