// Package priorityqueue describe an prioriry queue implement.
package priorityqueue

import (
	"container/heap"
	"fmt"

	"github.com/lsytj0413/ena/xerrors"
)

// Element for priority queue element
type Element[T any] struct {
	// Value for element
	Value T

	// priority of the element, make it private to avoid change from element
	priority int64

	// index of the element in the slice
	index int

	// pq is the refer of priority queue
	pq *priorityQueue[T]
}

// Priority return the priority value of element
func (e *Element[T]) Priority() int64 {
	return e.priority
}

// Index return the element index in slice
func (e *Element[T]) Index() int {
	return e.index
}

// String return the representation of element
func (e *Element[T]) String() string {
	return fmt.Sprintf("*Element{Value:%v,priority:%v,index:%v,pq:%v}", e.Value, e.priority, e.index, e.pq)
}

// PriorityQueue for priority queue trait
type PriorityQueue[T any] interface {
	// Add element to the PriorityQueue, it will return the element witch been added
	Add(v T, priority int64) *Element[T]

	// Peek return the lowest priority element
	Peek() *Element[T]

	// Pop return the lowest priority element and remove it
	Pop() *Element[T]

	// Remove will remove the element from the priority queue
	Remove(v *Element[T]) error

	// Update the element in the priority queue with the new priority
	Update(v *Element[T], priority int64) error

	// Size return the element size of queue
	Size() int
}

// heapi implement the heap.Interface
type heapi[T any] struct {
	pq *priorityQueue[T]
}

// Len return the size of slice, implement heap.Len
func (h *heapi[T]) Len() int {
	return len(h.pq.e)
}

// Less return the compare of slice element i and j, implement heap.Less
func (h *heapi[T]) Less(i int, j int) bool {
	return h.pq.e[i].priority < h.pq.e[j].priority
}

// Swap change the slice element i and j, implement heap.Swap
func (h *heapi[T]) Swap(i int, j int) {
	h.pq.e[i], h.pq.e[j] = h.pq.e[j], h.pq.e[i]
	h.pq.e[i].index, h.pq.e[j].index = h.pq.e[j].index, h.pq.e[i].index
}

// Push the value at the end of slice, implement heap.Push
func (h *heapi[T]) Push(x interface{}) {
	h.pq.e = append(h.pq.e, x.(*Element[T]))
}

// Pop the value at the last position of slice, implement heap.Pop
func (h *heapi[T]) Pop() interface{} {
	old := h.pq.e
	n := len(old)
	if n == 0 {
		return nil
	}

	// set element to nil for GC
	x := old[n-1]
	old[n-1] = nil
	h.pq.e = old[0 : n-1]
	return x
}

// priorityQueue is a implement by min heap, the 0th element is the lowest value
type priorityQueue[T any] struct {
	e []*Element[T]
	h *heapi[T]
}

// NewPriorityQueue construct a PriorityQueue
func NewPriorityQueue[T any](size int) PriorityQueue[T] {
	pq := &priorityQueue[T]{
		e: make([]*Element[T], 0, size),
	}
	pq.h = &heapi[T]{
		pq: pq,
	}
	return pq
}

// Add element to the PriorityQueue, it will return the element witch been added
func (pq *priorityQueue[T]) Add(x T, priority int64) *Element[T] {
	e := &Element[T]{
		Value:    x,
		priority: priority,
		index:    len(pq.e),
		pq:       pq,
	}
	heap.Push(pq.h, e)
	return e
}

// Peek return the lowest priority element
func (pq *priorityQueue[T]) Peek() *Element[T] {
	if len(pq.e) == 0 {
		return nil
	}

	return pq.e[0]
}

// Pop return the lowest priority element and remove it
func (pq *priorityQueue[T]) Pop() *Element[T] {
	if len(pq.e) == 0 {
		return nil
	}

	x := heap.Pop(pq.h)
	e := x.(*Element[T])
	e.index = -1
	e.pq = nil
	return e
}

var (
	// ErrMismatchQueue represent the element's queue mismatch current operation queue
	ErrMismatchQueue = xerrors.Errorf("MismatchQueue")

	// ErrOutOfIndex represent the element's index out of queue's length
	ErrOutOfIndex = xerrors.Errorf("OutOfIndex")

	// ErrMismatchPriority represent the element's priority mismatch
	ErrMismatchPriority = xerrors.Errorf("MismatchPriority")
)

// Remove will remove the element from the priority queue
func (pq *priorityQueue[T]) Remove(e *Element[T]) error {
	if e.pq != pq {
		return xerrors.Wrapf(ErrMismatchQueue, "element %p, current %p", e.pq, pq)
	}

	if e.index < 0 || e.index >= len(pq.e) {
		return xerrors.Wrapf(ErrOutOfIndex, "element index %d, length %d", e.index, len(pq.e))
	}
	if e.priority != pq.e[e.index].priority {
		return xerrors.Wrapf(ErrMismatchPriority, "element priority %d, current %d", e.priority, pq.e[e.index].priority)
	}

	heap.Remove(pq.h, e.index)
	e.index = -1
	e.pq = nil
	return nil
}

// Update the element in the priority queue with the new priority
func (pq *priorityQueue[T]) Update(e *Element[T], priority int64) error {
	if e.pq != pq {
		return xerrors.Wrapf(ErrMismatchQueue, "element %p, current %p", e.pq, pq)
	}
	if e.index < 0 || e.index >= len(pq.e) {
		return xerrors.Wrapf(ErrOutOfIndex, "element index %d, length %d", e.index, len(pq.e))
	}
	if e.priority == priority {
		// the priority doesn't change, just return nil as updated
		return nil
	}

	e.priority = priority
	heap.Fix(pq.h, e.index)
	return nil
}

func (pq *priorityQueue[T]) Size() int {
	return len(pq.e)
}
