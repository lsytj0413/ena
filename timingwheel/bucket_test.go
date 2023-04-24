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

package timingwheel

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBucketSetExpiration(t *testing.T) {
	g := NewWithT(t)
	b := newBucket()

	g.Expect(b.expiration).To(Equal(int64(-1)))

	g.Expect(b.SetExpiration(int64(-1))).To(BeFalse())
	g.Expect(b.expiration).To(Equal(int64(-1)))

	g.Expect(b.SetExpiration(int64(1))).To(BeTrue())
	g.Expect(b.expiration).To(Equal(int64(1)))
}

func TestBucketAdd(t *testing.T) {
	g := NewWithT(t)
	b := newBucket()

	ti := &timerTask{}
	b.Add(ti)
	g.Expect(b).To(Equal(ti.b))
	g.Expect(ti.e).ToNot(BeNil())
}

func TestBucketRemove(t *testing.T) {
	t.Run("failed", func(t *testing.T) {
		g := NewWithT(t)
		b := newBucket()

		ti := &timerTask{}
		b.Add(ti)

		ti.b = nil
		g.Expect(b.remove(ti)).To(BeFalse())
		g.Expect(ti.b).To(BeNil())
		g.Expect(ti.e).ToNot(BeNil())
	})

	t.Run("ok", func(t *testing.T) {
		g := NewWithT(t)
		b := newBucket()

		ti := &timerTask{}
		b.Add(ti)

		g.Expect(b.remove(ti)).To(BeTrue())
		g.Expect(ti.b).To(BeNil())
		g.Expect(ti.e).To(BeNil())
	})
}

func TestBucketFlush(t *testing.T) {
	g := NewWithT(t)
	b := newBucket()

	timers := []*timerTask{
		{},
		{},
	}
	for _, ti := range timers {
		b.Add(ti)
	}

	gotTimers := make([]*timerTask, 0, len(timers))
	reinsert := func(t *timerTask) {
		gotTimers = append(gotTimers, t)
	}

	b.Flush(reinsert)

	for _, ti := range timers {
		g.Expect(ti.b).To(BeNil())
		g.Expect(ti.e).To(BeNil())
	}

	for i := range timers {
		g.Expect(timers[i]).To(Equal(gotTimers[i]))
	}
}
