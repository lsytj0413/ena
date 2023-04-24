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

	"github.com/lsytj0413/ena/xerrors"
)

type testStopWheel struct {
	stopFuncFn func(*timerTask) (bool, error)
}

func (t *testStopWheel) StopFunc(tt *timerTask) (bool, error) {
	return t.stopFuncFn(tt)
}

func TestTimerTaskStop(t *testing.T) {
	t.Run("stopped", func(t *testing.T) {
		g := NewWithT(t)
		tt := &timerTask{}
		tt.stopped = 1

		v, err := tt.Stop()
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v).To(BeTrue())
	})

	t.Run("stop_failed", func(t *testing.T) {
		g := NewWithT(t)
		tt := &timerTask{
			w: &testStopWheel{
				stopFuncFn: func(tt *timerTask) (bool, error) {
					return false, xerrors.ErrContinue
				},
			},
		}

		v, err := tt.Stop()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`Continue`))
		g.Expect(v).To(BeFalse())
	})

	t.Run("ok", func(t *testing.T) {
		g := NewWithT(t)
		tt := &timerTask{
			w: &testStopWheel{
				stopFuncFn: func(tt *timerTask) (bool, error) {
					return true, nil
				},
			},
		}

		v, err := tt.Stop()
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v).To(BeTrue())
	})
}
