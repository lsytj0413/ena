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
