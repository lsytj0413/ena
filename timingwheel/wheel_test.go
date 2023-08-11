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
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	delayqueuemocks "github.com/lsytj0413/ena/delayqueue/mocks"
)

func TestNewWheel(t *testing.T) {
	g := NewWithT(t)
	w := newWheel(3, 20, 4)

	g.Expect(w.tick).To(Equal(int64(3)))
	g.Expect(w.wheelSize).To(Equal(int64(20)))
	g.Expect(w.interval).To(Equal(int64(60)))
	g.Expect(w.currentTime).To(Equal(int64(3)))
	g.Expect(w.buckets).To(HaveLen(int(20)))

	for _, b := range w.buckets {
		g.Expect(b).ToNot(BeNil())
	}
}

func TestWheelAdd(t *testing.T) {
	t.Run("false", func(t *testing.T) {
		w := newWheel(3, 20, 4)

		for i := int64(0); i < w.currentTime+w.tick; i++ {
			desp := fmt.Sprintf("exp %v", i)
			t.Run(desp, func(t *testing.T) {
				g := NewWithT(t)

				g.Expect(w.add(&timerTask{
					expiration: int64(i),
				}, nil)).To(BeFalse())
			})
		}
	})

	t.Run("ok", func(t *testing.T) {
		type testCase struct {
			desp             string
			t                *timerTask
			bucketIndex      int
			bucketExpiration int64
			isOfferCalled    bool
			currentTime      int64
		}
		testCases := []testCase{
			{
				desp: "task expiration 7, 2th bucket",
				t: &timerTask{
					expiration: 7,
				},
				bucketIndex:      2,
				bucketExpiration: 6,
				isOfferCalled:    true,
			},
			{
				desp: "task expiration 17, 5th bucket",
				t: &timerTask{
					expiration: 17,
				},
				bucketIndex:      5,
				bucketExpiration: 15,
				isOfferCalled:    true,
			},
			{
				desp: "task expiration 62, 0th bucket",
				t: &timerTask{
					expiration: 62,
				},
				bucketIndex:      0,
				bucketExpiration: 60,
				isOfferCalled:    true,
			},
			{
				desp: "task expiration 8, 2th bucket",
				t: &timerTask{
					expiration: 8,
				},
				bucketIndex:      2,
				bucketExpiration: 6,
				isOfferCalled:    false,
			},
			{
				desp:        "task expiration 16, 5th bucket",
				currentTime: 6,
				t: &timerTask{
					expiration: 16,
				},
				bucketIndex:      5,
				bucketExpiration: 15,
				isOfferCalled:    false,
			},
		}

		w := newWheel(3, 20, 4)
		mockCtrl := gomock.NewController(t)
		dq := delayqueuemocks.NewMockDelayQueue[*bucket](mockCtrl)

		for _, tc := range testCases {
			t.Run(tc.desp, func(t *testing.T) {
				g := NewWithT(t)

				if tc.currentTime > 0 {
					w.advanceClock(tc.currentTime)
				}
				if tc.isOfferCalled {
					dq.EXPECT().Offer(gomock.Any(), gomock.Any()).Times(1)
				} else {
					dq.EXPECT().Offer(gomock.Any(), gomock.Any()).Times(0)
				}

				g.Expect(w.add(tc.t, dq)).To(BeTrue())
				g.Expect(w.buckets[tc.bucketIndex]).To(Equal(tc.t.b))
				g.Expect(w.buckets[tc.bucketIndex].Expiration()).To(Equal(tc.bucketExpiration))
			})
		}
	})

	t.Run("ok_overflow_wheel", func(t *testing.T) {
		g := NewWithT(t)
		w := newWheel(3, 20, 4)
		mockCtrl := gomock.NewController(t)
		dq := delayqueuemocks.NewMockDelayQueue[*bucket](mockCtrl)

		g.Expect(w.overflowWheel).To(BeNil())

		// the overflowwheel will cover the uplayerwheel range. at this EX:
		// 1. uplayer: tick(3), current(3), wheelsize(20), interval(60)
		// 2. uplayer: range(3-63)
		// 3. overflow layer:
		//   a. tick(60): uplayer interval
		//   b. current(0): truncate(uplayer current, tick)
		//   c. wheelsize(20): uplayer wheelsize
		//   d. interval(1200): tick*wheelsize
		// 4. overflow layer: range(0-1200)
		dq.EXPECT().Offer(gomock.Any(), gomock.Any()).Times(1)

		g.Expect(w.add(&timerTask{
			expiration: 66,
		}, dq)).To(BeTrue())
		g.Expect(w.overflowWheel).ToNot(BeNil())
		g.Expect(w.overflowWheel.buckets[1]).ToNot(BeNil())
		g.Expect(w.overflowWheel.buckets[1].Expiration()).To(Equal(int64(60)))
	})
}

func TestWheelAddOrRun(t *testing.T) {
	t.Run("run immediate after task", func(t *testing.T) {
		g := NewWithT(t)
		w := newWheel(3, 20, 4)

		mockCtrl := gomock.NewController(t)
		dq := delayqueuemocks.NewMockDelayQueue[*bucket](mockCtrl)
		dq.EXPECT().Offer(gomock.Any(), gomock.Any()).Times(0)

		v := 0
		tt := &timerTask{
			expiration: 0,
			f: func(time.Time) {
				v = 1
			},
			t: taskAfter,
		}
		w.addOrRun(tt, dq)
		g.Expect(v).To(Equal(1))
		g.Expect(tt.b).To(BeNil())
	})

	t.Run("run immediate nonstoped tick task", func(t *testing.T) {
		g := NewWithT(t)
		w := newWheel(3, 20, 4)

		v := 0
		tt := &timerTask{
			expiration: 0,
			d:          time.Duration(7),
			stopped:    0,
			f: func(time.Time) {
				v = 1
			},
			t: taskTick,
		}

		mockCtrl := gomock.NewController(t)
		dq := delayqueuemocks.NewMockDelayQueue[*bucket](mockCtrl)
		dq.EXPECT().Offer(gomock.Any(), gomock.Any()).Times(1)

		w.addOrRun(tt, dq)
		g.Expect(v).To(Equal(1))
		g.Expect(tt.b).ToNot(BeNil())
	})

	t.Run("run immediate stopped tick task", func(t *testing.T) {
		g := NewWithT(t)
		w := newWheel(3, 20, 4)

		v := 0
		tt := &timerTask{
			expiration: 0,
			d:          time.Duration(7),
			stopped:    1,
			f: func(time.Time) {
				v = 1
			},
			t: taskTick,
		}

		mockCtrl := gomock.NewController(t)
		dq := delayqueuemocks.NewMockDelayQueue[*bucket](mockCtrl)
		dq.EXPECT().Offer(gomock.Any(), gomock.Any()).Times(0)

		w.addOrRun(tt, dq)
		g.Expect(v).To(Equal(1))
		g.Expect(tt.b).To(BeNil())
	})
}

func TestWheelAdvanceClock(t *testing.T) {
	g := NewWithT(t)
	w := newWheel(3, 20, 4)
	w.overflowWheel = newWheel(0, 0, 0)

	g.Expect(w.overflowWheel.currentTime).To(Equal(int64(0)))

	exp := int64(100)
	w.advanceClock(exp)
	exp = int64(99)
	g.Expect(w.currentTime).To(Equal(exp))
	g.Expect(w.overflowWheel.currentTime).To(Equal(exp))
}
