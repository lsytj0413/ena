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

	. "github.com/onsi/gomega"
)

func TestWaitRegister(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := NewWithT(t)
		w := New().(*defWait)

		for i := uint64(0); i < uint64(100); i++ {
			_, err := w.Register(strconv.FormatUint(i, 10))
			g.Expect(err).ToNot(HaveOccurred())
		}

		for i := uint64(0); i < uint64(100); i++ {
			ok := w.IsRegistered(strconv.FormatUint(i, 10))
			g.Expect(ok).To(BeTrue())
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		g := NewWithT(t)
		w := New().(*defWait)

		indexs := []uint64{0, 99, 100, 234234, 34535435345345, ^uint64(0)}
		for _, i := range indexs {
			_, err := w.Register(strconv.FormatUint(i, 10))
			g.Expect(err).ToNot(HaveOccurred())
		}

		for _, i := range indexs {
			_, err := w.Register(strconv.FormatUint(i, 10))
			g.Expect(err).To(HaveOccurred())
			g.Expect(err.Error()).To(MatchRegexp(`DuplicateObject`))
		}
	})
}

func TestWaitTrigger(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		g := NewWithT(t)
		w := New().(*defWait)

		id := "98989"
		ch, err := w.Register(id)
		g.Expect(err).ToNot(HaveOccurred())

		err = w.Trigger(id, id)
		g.Expect(err).ToNot(HaveOccurred())

		// check value
		v := <-ch
		v0, ok := v.(string)
		g.Expect(ok).To(BeTrue())
		g.Expect(v0).To(Equal(id))

		// check register
		ok = w.IsRegistered(id)
		g.Expect(ok).To(BeFalse())
	})

	t.Run("not_register", func(t *testing.T) {
		g := NewWithT(t)
		w := New().(*defWait)

		id := "98989"
		err := w.Trigger(id, id)
		g.Expect(err).To(HaveOccurred())

		g.Expect(err.Error()).To(MatchRegexp(`NotFound`))
	})

	t.Run("timeout", func(t *testing.T) {
		g := NewWithT(t)
		w := New().(*defWait)

		id := "98989"
		w.Register(id)

		ch := w.m[id]
		ch <- nil

		err := w.Trigger(id, id)
		g.Expect(err).To(HaveOccurred())

		g.Expect(err.Error()).To(MatchRegexp(`timeout`))
	})
}
