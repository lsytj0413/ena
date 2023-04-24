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

// Package timingwheel implementation of Hierarchical Timing Wheels.
package timingwheel

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/lsytj0413/ena"
	"github.com/lsytj0413/ena/delayqueue"
	"github.com/lsytj0413/ena/wait"
)

type option struct {
	// Tick is the tick duration with wheel
	Tick time.Duration

	// WheelSize is the size of buckets
	WheelSize int64
}

// Validate check the option
func (o *option) Validate() error {
	return nil
}

// Option is some configuration that modifies options for a wheel.
type Option interface {
	Apply(*option)
}

// WithTickDuration set the Tick field
type WithTickDuration time.Duration

// Apply applies this configuration to the given option
func (w WithTickDuration) Apply(opt *option) {
	opt.Tick = time.Duration(w)
}

// WithSize set the WheelSize field
type WithSize int64

// Apply applies this configuration to the given option
func (w WithSize) Apply(opt *option) {
	opt.WheelSize = int64(w)
}

// NewTimingWheel creates an instance of TimingWheel with the given tick and wheelSize.
func NewTimingWheel(opts ...Option) (TimingWheel, error) {
	options := &option{
		Tick:      time.Second,
		WheelSize: 64,
	}
	for _, opt := range opts {
		opt.Apply(options)
	}

	tickMs := int64(options.Tick / time.Millisecond)
	if tickMs <= 0 {
		return nil, ErrInvalidTickValue
	}
	if options.WheelSize <= 0 {
		return nil, ErrInvalidWheelSize
	}

	startMs := timeToMs(time.Now())
	t := newWheel(tickMs, options.WheelSize, startMs)

	tw := &timingWheel{
		dq:  delayqueue.New[*bucket](int(options.WheelSize)),
		w:   t,
		wt:  wait.New(),
		wch: make(chan event, int(options.WheelSize)*100), // the channel is bufferd, could change to unbufferd?
	}

	tw.ctx, tw.cancel = context.WithCancel(context.Background())
	return tw, nil
}

// timingWheel is an implemention of TimingWheel
// TODO(lsytj0413): disable add/tick/stopfunc when not running?
type timingWheel struct {
	// the first layer wheel
	w *wheel

	// the Wait to get response when call AfterFunc
	wt wait.Wait

	// the Wait register id, incr
	wid uint64

	// wch is the channel which TimerTask putin when call AfterFunc
	wch chan event

	// dq is the queue of bucket expiration
	dq delayqueue.DelayQueue[*bucket]

	// wg for wait sub goroutine
	wg ena.WaitGroupWrapper

	// ctx to cancel sub goroutine
	ctx    context.Context
	cancel func()
}

// event is the representation of value in the wch
type event struct {
	// the identify of the event
	Type eventType

	t *timerTask
}

// Start will start the timingwheel, and process the tasks
// nolint
func (tw *timingWheel) Start() {
	tw.wg.Wrap(func() {
		tw.dq.Poll(tw.ctx)
	})

	addOrRun := func(t *timerTask) {
		tw.w.addOrRun(t, tw.dq)
	}

	tw.wg.Wrap(func() {
		for {
			select {
			case b := <-tw.dq.Chan():
				tw.w.advanceClock(b.Expiration())

				b.Flush(addOrRun)
			case e := <-tw.wch:
				switch e.Type {
				case eventAddNew:
					// an timer task is add from AfterFunc/StopFunc
					addOrRun(e.t)
					_ = tw.wt.Trigger(strconv.FormatUint(e.t.id, 10), e.t)
				case eventDelete:
					stopped := false
					switch atomic.LoadUint32(&e.t.stopped) {
					case 1:
						stopped = true
					default:
						if e.t.b != nil {
							stopped = e.t.b.remove(e.t)
						}
						if stopped {
							atomic.StoreUint32(&e.t.stopped, 1)
						}
					}
					_ = tw.wt.Trigger(strconv.FormatUint(e.t.id, 10), stopped)
				}
			case <-tw.ctx.Done():
				return
			}
		}
	})
}

// TODO(lsytj0413): clear all timer after stop?
func (tw *timingWheel) Stop() {
	tw.cancel()
	tw.wg.Wait()
}

func (tw *timingWheel) AfterFunc(d time.Duration, f Handler) (TimerTask, error) {
	return tw.addFunc(d, f, taskAfter)
}

func (tw *timingWheel) TickFunc(d time.Duration, f Handler) (TimerTask, error) {
	v := d / (time.Duration(tw.w.tick) * time.Millisecond)
	if v <= 0 {
		return nil, ErrInvalidTickFuncDurationValue
	}

	return tw.addFunc(d, f, taskTick)
}

func (tw *timingWheel) StopFunc(t *timerTask) (bool, error) {
	outch, err := tw.enquenTask(t, eventDelete)
	if err != nil {
		return false, err
	}

	v := <-outch
	return v.(bool), nil
}

func (tw *timingWheel) addFunc(d time.Duration, f Handler, eType timerTaskType) (TimerTask, error) {
	t := &timerTask{
		d:          d,
		expiration: timeToMs(time.Now().Add(d)),
		t:          eType,
		f:          f,
		id:         atomic.AddUint64(&tw.wid, 1),
		w:          tw,
	}

	outch, err := tw.enquenTask(t, eventAddNew)
	if err != nil {
		return nil, err
	}

	v := <-outch
	return v.(*timerTask), nil
}

func (tw *timingWheel) enquenTask(t *timerTask, eType eventType) (<-chan interface{}, error) {
	outch, err := tw.wt.Register(strconv.FormatUint(t.id, 10))
	if err != nil {
		return nil, err
	}

	tw.wch <- event{
		Type: eType,
		t:    t,
	}

	return outch, nil
}
