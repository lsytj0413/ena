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

package wait

import (
	"strconv"
	"testing"
	"time"
)

func Benchmark_Wait(b *testing.B) {
	w := New()

	for i := 0; i < b.N; i++ {
		id := strconv.FormatUint(uint64(i), 10)
		w.Register(id)
		w.Trigger(id, nil)
	}
}

func Benchmark_Wait_Parallel(b *testing.B) {
	w := New()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := uint64(time.Now().UnixNano())
			id := strconv.FormatUint(uint64(i), 10)

			w.Register(id)
			w.Trigger(id, nil)
		}
	})
}
