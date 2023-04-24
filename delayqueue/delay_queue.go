// Package delayqueue describe an delayqueue implemention.
package delayqueue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lsytj0413/ena/priorityqueue"
)

// DelayQueue is an blocking queue of *Delay* elements, the element
// can only been taken when its delay has expired. The head of the queue
// is the element whose delay expired most recent in the queue.
type DelayQueue[T any] interface {
	// Offer insert the element into the current DelayQueue,
	// if the expiration is blow the current min expiration, the item will
	// been fired first.
	Offer(elem T, expireation int64)

	// Poll starts an infinite loop, it will continually waits for an element to
	// been fired, and send the element to the output Chan.
	Poll(ctx context.Context)

	// Chan return the output chan, when the element is fired the element
	// will send to the channel.
	Chan() <-chan T

	// Size return the element count in the queue
	Size() int
}

// delayQueue implement the DelayQueue interface
type delayQueue[T any] struct {
	// C is the output channel, when element is fired it will send into this channel
	C chan T

	// wakeupC is the inner channel for wakeup the sleep, when new item is add into the
	// queue, the min expired maybe change, and then the channel will be readable
	wakeupC chan struct{}

	// T is the timer to provide the current ms
	T Timer

	// pq is the priorityqueue of expiration
	pq priorityqueue.PriorityQueue[T]

	// protect the add/remove/update operation in the PriorityQueue
	// TODO(yangsonglin): implement the goroutine-safe PriorityQueue
	mu sync.Mutex

	// sleeping is the sleeping state of delayqueue, if the queue is waiting for fired, the value will be 1
	sleeping int32

	// for unittest
	pollFn func(ctx context.Context, q *delayQueue[T]) bool
}

// New construct a DelayQueue with the initial size
func New[T any](size int) DelayQueue[T] {
	return NewWithTimer[T](size, defaultTimer)
}

// NewWithTimer construct a DelayQueue with the initial size and Timer
func NewWithTimer[T any](size int, t Timer) DelayQueue[T] {
	return &delayQueue[T]{
		C:       make(chan T),
		wakeupC: make(chan struct{}),
		T:       t,
		pq:      priorityqueue.NewPriorityQueue[T](1),
		pollFn:  pollImpl[T],
	}
}

// TODO(yangsonglin): is't too difficult to deal with the sleeping, so change to the worker model?
// Offer implement the DelayQueue.Offer
func (q *delayQueue[T]) Offer(element T, expireation int64) {
	_push := func() (*priorityqueue.Element[T], int) {
		q.mu.Lock()
		defer q.mu.Unlock()

		e := q.pq.Add(element, expireation)
		return e, e.Index()
	}
	_, index := _push()

	// there is no concurrent protection, EX:
	// 1. goroutine1 add element with expireation 100
	// 2. goroutine2 add element with expireation 50
	// 3. the both goroutine get the element index 0
	// 4. goroutine2 cas the sleeping state to 0, and send the wakeup signal
	// 5. pool wakeup and update the fired point, cas the sleeping state to 1
	// 6. goroutine1 cas the sleeping state to 0, and send the wakeup signal
	// 7. pool wakeup and update the fired point, cas the sleeping state to 1
	// because the pool always update the fired point to the min expireation, so there is no problem(always update to 50)
	if index == 0 {
		// the element is the first element(with the earliest expireation), we
		// need week up the Pool loop to update the fired point
		if atomic.CompareAndSwapInt32(&q.sleeping, 1, 0) {
			// if we change the sleeping state from sleep to weekup success, send the signal to wakepupC
			q.wakeupC <- struct{}{}
		}
	}
}

// Poll implement the DelayQueue.Pool
func (q *delayQueue[T]) Poll(ctx context.Context) {
	defer func() {
		// reset the state to wakeup
		atomic.StoreInt32(&q.sleeping, 0)
	}()

	// an infinite loop
	// 1. wakeup at the min expiration
	// 2. send to the C
	for q.pollFn(ctx, q) {
	}
}

// the inner implement of poll, split from Poll for test
// return true if been wakeup or fired, false to shutdown the loop
// nolint
func pollImpl[T any](ctx context.Context, q *delayQueue[T]) bool {
	n := q.T.Now()

	q.mu.Lock()
	item := q.pq.Peek()
	if item == nil || item.Priority() > n {
		// No item left, change the sleeping state to 1
		atomic.StoreInt32(&q.sleeping, 1)
	}
	q.mu.Unlock()

	// we have got the min expiration item, it maybe nil for empty pq
	if item == nil {
		// wait for wakeup (new item Offer into the queue)
		select {
		case <-ctx.Done():
			return false
		case <-q.wakeupC:
			return true
		}
	}

	// have item, wait for the fired point
	delta := item.Priority() - n
	if delta <= 0 {
		// the item need fired, send the value to the output channel
		select {
		// TODO(yangsonglin): change to executor
		case q.C <- item.Value:
			// the element is fired
			q.mu.Lock()
			_ = q.pq.Remove(item)
			q.mu.Unlock()
			return true
		case <-ctx.Done():
			return false
		}
	}

	// the item is pending, wait for fired or new min element add
	select {
	case <-q.wakeupC:
		return true
	case <-time.After(time.Duration(delta) * time.Millisecond):
		// we doesn't fired the item at there, go to next loop and the item will been fired because delta <= 0
		if atomic.SwapInt32(&q.sleeping, 0) == 0 {
			// if the old state is wakeup, the maybe an signal in wakeupC,
			// so we drain it the unblock the caller
			select {
			case <-q.wakeupC:
			default:
			}
		}
		return true
	case <-ctx.Done():
		return false
	}
}

// Chan implement the DelayQueue.Chan
func (q *delayQueue[T]) Chan() <-chan T {
	return q.C
}

// Size implement the DelayQueue.Size
func (q *delayQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.pq.Size()
}
