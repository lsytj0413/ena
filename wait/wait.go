// Package wait provides utility functions for polling, listening using Go channel
package wait

import (
	"sync"

	"github.com/lsytj0413/ena/xerrors"
)

// Wait is an interface that provices the ability to wait and trigger events that
// are associated with ID
type Wait interface {
	// Register waits returns a chan that waits on the given ID.
	// The chan will be triggered when Trigger is called with the same ID
	Register(id string) (<-chan interface{}, error)

	// Trigger triggers the waiting chans with the given ID
	Trigger(id string, x interface{}) error

	// IsRegisterd returns where the id is been registered
	IsRegistered(id string) bool
}

type defWait struct {
	sync.RWMutex
	m map[string]chan interface{}
}

const (
	defaultMapSize = 1000
)

// New creates a Wait object
func New() Wait {
	return &defWait{
		m: make(map[string]chan interface{}, defaultMapSize),
	}
}

func (w *defWait) Register(id string) (<-chan interface{}, error) {
	w.Lock()
	defer w.Unlock()

	if _, ok := w.m[id]; !ok {
		c := make(chan interface{}, 1)
		w.m[id] = c
		return c, nil
	}

	return nil, xerrors.WrapDuplicate("Wait.Register: id %s", id)
}

func (w *defWait) Trigger(id string, x interface{}) error {
	f := func() (chan interface{}, bool) {
		w.Lock()
		defer w.Unlock()

		c, ok := w.m[id]
		delete(w.m, id)
		return c, ok
	}

	c, ok := f()
	if !ok {
		return xerrors.WrapNotFound("Wait.Trigger: id %s", id)
	}

	select {
	case c <- x:
		close(c)
		return nil
	default:
		return xerrors.Errorf("WaitTrigger: id %s timeout", id)
	}
}

func (w *defWait) IsRegistered(id string) bool {
	w.RLock()
	defer w.RUnlock()

	_, ok := w.m[id]
	return ok
}
