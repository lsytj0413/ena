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

package delayqueue

import (
	"context"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestDelayQueueOffer(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer func() {
			t.Logf("Poll done")
			wg.Done()
		}()

		dq.Poll(ctx)
	}()
	go func() {
		defer func() {
			t.Logf("Receive done")
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case v := <-dq.Chan():
				t.Logf("Receive %v", v)
			}
		}
	}()

	time.Sleep(10 * time.Millisecond)
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))

	n := defaultTimer.Now()
	n += 1000
	n += 10
	t.Logf("Offer first element: priority %v", n)
	dq.Offer(int(n), n)
	g.Expect(dq.Size()).To(Equal(1))

	// TODO(yangsonglin): difficult to test, we expect to see 0
	v := atomic.LoadInt32(&dq.sleeping)
	if v != 0 && v != 1 {
		g.Ω(func() bool {
			return v == 0 || v == 1
		}).To(BeTrue(), "sleeping[%v] must be 0 or 1", v)
	}

	// wait be sleeping
	time.Sleep(10 * time.Millisecond)
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))

	// won't been wakeup
	n += 20
	t.Logf("Offer second element: priority %v", n)
	dq.Offer(int(n), n)
	g.Expect(dq.Size()).To(Equal(2))
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))

	// been wakeup
	n -= 40
	t.Logf("Offer third element: priority %v", n)
	dq.Offer(int(n), n)
	g.Expect(dq.Size()).To(Equal(3))
	// TODO(yangsonglin): difficult to test, we expect to see 0
	v = atomic.LoadInt32(&dq.sleeping)
	if v != 0 && v != 1 {
		g.Ω(func() bool {
			return v == 0 || v == 1
		}).To(BeTrue(), "sleeping[%v] must be 0 or 1", v)
	}

	// wait be sleeping
	time.Sleep(10 * time.Millisecond)
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))

	cancel()
	wg.Wait()
}

func TestDelayQueuePoll(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	dq.pollFn = func(ctx context.Context, q *delayQueue[int]) bool {
		return false
	}

	dq.Poll(context.Background())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(0)))
}

func TestDelayQueuePollImplNullItem(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	// cancel context
	n := defaultTimer.Now()
	r := pollImpl(ctx, dq)
	g.Expect(r).To(BeFalse())
	g.Ω(func() bool {
		return defaultTimer.Now()-n <= 2
	}()).To(BeTrue())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))

	// been wakeup
	go func() {
		dq.wakeupC <- struct{}{}
	}()
	r = pollImpl(context.Background(), dq)
	g.Expect(r).To(BeTrue())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))
}

func TestDelayQueuePollAfterDelayItem(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	// test the item should been fired
	dq.Offer(1, 0)
	g.Expect(dq.Size()).To(Equal(1))

	// fired
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		v := <-dq.Chan()
		g.Expect(v).To(Equal(1))
	}()

	r := pollImpl(context.Background(), dq)
	g.Expect(r).To(BeTrue())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(0)))
	g.Expect(dq.Size()).To(Equal(0))

	wg.Wait()

	// been wakeup
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	go func(wakeupC chan struct{}) {
		wakeupC <- struct{}{}
	}(dq.wakeupC)

	dq.Offer(1, 0)
	g.Expect(dq.Size()).To(Equal(1))

	r = pollImpl(ctx, dq)
	g.Expect(r).To(BeFalse())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(0)))
	g.Expect(dq.Size()).To(Equal(1))
}

func TestDelayQueuePollBeforeDelayItemFired(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	// test the item shouldn't been fired
	n := defaultTimer.Now() + 1000
	dq.Offer(1, n)
	g.Expect(dq.Size()).To(Equal(1))

	// wait been fired
	r := pollImpl(context.Background(), dq)
	g.Expect(r).To(BeTrue())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(0)))
	g.Expect(dq.Size()).To(Equal(1))
}

func TestDelayQueueBeforeDelayItemWakeup(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	// test the item shouldn't been fired
	n := defaultTimer.Now() + 1000

	// been wakeup
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	go func() {
		dq.wakeupC <- struct{}{}
	}()
	dq.Offer(1, n)
	g.Expect(dq.Size()).To(Equal(1))

	r := pollImpl(ctx, dq)
	g.Expect(r).To(BeTrue())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))
	g.Expect(dq.Size()).To(Equal(1))
}

func TestDelayQueuePollBeforeDelayItemCancel(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	// test the item shouldn't been fired
	n := defaultTimer.Now() + 1000

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dq.Offer(1, n)
	g.Expect(dq.Size()).To(Equal(1))

	r := pollImpl(ctx, dq)
	g.Expect(r).To(BeFalse())
	g.Expect(atomic.LoadInt32(&dq.sleeping)).To(Equal(int32(1)))
	g.Expect(dq.Size()).To(Equal(1))
}

func TestDelayQueueSize(t *testing.T) {
	g := NewWithT(t)
	dq := New[int](1).(*delayQueue[int])

	ctx, cancel := context.WithCancel(context.Background())

	type testCase struct {
		value       int
		expireation int64
	}
	delay := int64(2000) // 2000ms, 2s
	testCases := []testCase{
		{
			value:       1,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
		{
			value:       2,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
		{
			value:       3,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
		{
			value:       4,
			expireation: defaultTimer.Now() + int64(rand.Intn(int(delay))),
		},
	}
	sortTestCases := make([]testCase, len(testCases))
	copy(sortTestCases, testCases)
	sort.Slice(sortTestCases, func(i int, j int) bool {
		return sortTestCases[i].expireation < sortTestCases[j].expireation
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer func() {
			t.Logf("Poll done")
			wg.Done()
		}()

		dq.Poll(ctx)
	}()

	go func() {
		defer func() {
			t.Logf("Receive done")
			wg.Done()
		}()

		i := 0
		for v := range dq.Chan() {
			n := defaultTimer.Now()
			diff := n - sortTestCases[i].expireation
			g.Ω(func() bool {
				return diff < 0 || diff > 10
			}()).To(BeFalse(), "ExpireationPrecesionFailed", "Receive at[%v], estimate[%v]", n, sortTestCases[i].expireation)

			t.Logf("Receive %v", v)
			g.Expect(v).To(Equal(sortTestCases[i].value))
			i++
			if i >= len(sortTestCases) {
				cancel()
				return
			}
		}
	}()

	for i, tc := range testCases {
		dq.Offer(tc.value, tc.expireation)
		g.Expect(dq.Size()).To(Equal(i + 1))
	}

	wg.Wait()

	g.Expect(dq.Size()).To(Equal(0))
}
