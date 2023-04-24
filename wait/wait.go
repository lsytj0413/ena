// Copyright (c) 2023 The Songlin Yang Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
