package ena

import (
	"sync"
)

// WaitGroupWrapper is a help struct for sync.WaitGroup
type WaitGroupWrapper struct {
	NoCopy
	wg sync.WaitGroup
}

// Wrap will handle sync.WaitGroup Add and Done
func (w *WaitGroupWrapper) Wrap(cb func(), cbs ...func()) {
	cbs = append(cbs, cb)
	w.wg.Add(len(cbs))

	for _, cb := range cbs {
		go func(f func()) {
			defer w.wg.Done()
			f()
		}(cb)
	}
}

// Wait will block until all cb finished
func (w *WaitGroupWrapper) Wait() {
	w.wg.Wait()
}
