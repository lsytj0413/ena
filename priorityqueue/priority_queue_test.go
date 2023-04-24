package priorityqueue

import (
	"container/heap"
	"math/rand"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPriorityQueueAddOk(t *testing.T) {
	pq := NewPriorityQueue[int](0).(*priorityQueue[int])
	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority 1",
			value:       500,
			priority:    1,
			index:       1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			e := pq.Add(tc.value, tc.priority)
			g.Expect(tc.value).To(Equal(e.Value))
			g.Expect(tc.index).To(Equal(e.index))
			g.Expect(e.Priority()).To(Equal(tc.priority))
		})
	}

	testCases2 := []testCase{
		{
			description: "index 0",
			index:       0,
			priority:    0,
			value:       300,
		},
		{
			description: "index 1",
			index:       1,
			priority:    1,
			value:       500,
		},
		{
			description: "index 2",
			index:       2,
			priority:    1,
			value:       100,
		},
		{
			description: "index 3",
			index:       3,
			priority:    100,
			value:       400,
		},
		{
			description: "index 4",
			index:       4,
			priority:    2,
			value:       200,
		},
	}
	for i, e := range pq.e {
		tc := testCases2[i]
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(e.index).To(Equal(i))
			g.Expect(e.index).To(Equal(tc.index))
			g.Expect(e.priority).To(Equal(tc.priority))
			g.Expect(e.Value).To(Equal(tc.value))
		})
	}
}

func TestPriorityQueuePeek(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])
		g := NewWithT(t)

		g.Expect(pq.Peek()).To(BeNil())
	})

	t.Run("ok", func(t *testing.T) {
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		type testCase struct {
			description string
			value       int
			priority    int64
			index       int

			targetValue    int
			targetPriority int64
		}
		testCases := []testCase{
			{
				description:    "add priority 1",
				value:          100,
				priority:       1,
				index:          0,
				targetValue:    100,
				targetPriority: 1,
			},
			{
				description:    "add priority 2",
				value:          200,
				priority:       2,
				index:          1,
				targetValue:    100,
				targetPriority: 1,
			},
			{
				description:    "add priority 0",
				value:          300,
				priority:       0,
				index:          0,
				targetValue:    300,
				targetPriority: 0,
			},
			{
				description:    "add priority 100",
				value:          400,
				priority:       100,
				index:          3,
				targetValue:    300,
				targetPriority: 0,
			},
			{
				description:    "add priority -1",
				value:          500,
				priority:       -1,
				index:          0,
				targetValue:    500,
				targetPriority: -1,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				g := NewWithT(t)

				e := pq.Add(tc.value, tc.priority)
				g.Expect(e.Value).To(Equal(tc.value))
				g.Expect(e.index).To(Equal(tc.index))
				g.Expect(e.priority).To(Equal(tc.priority))

				e = pq.Peek()
				g.Expect(e).ToNot(BeNil())
				g.Expect(e.Value).To(Equal(tc.targetValue))
				g.Expect(e.Priority()).To(Equal(tc.targetPriority))
			})
		}
	})
}

func TestPriorityQueuePop(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])
		g := NewWithT(t)

		g.Expect(pq.Pop()).To(BeNil())
	})

	t.Run("ok", func(t *testing.T) {
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		type testCase struct {
			description string
			value       int
			priority    int64
			index       int
		}
		testCases := []testCase{
			{
				description: "add priority 1",
				value:       100,
				priority:    1,
				index:       0,
			},
			{
				description: "add priority 2",
				value:       200,
				priority:    2,
				index:       1,
			},
			{
				description: "add priority 0",
				value:       300,
				priority:    0,
				index:       0,
			},
			{
				description: "add priority 100",
				value:       400,
				priority:    100,
				index:       3,
			},
			{
				description: "add priority -1",
				value:       500,
				priority:    -1,
				index:       0,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				g := NewWithT(t)

				e := pq.Add(tc.value, tc.priority)
				g.Expect(e.Value).To(Equal(tc.value))
				g.Expect(e.index).To(Equal(tc.index))
				g.Expect(e.Priority()).To(Equal(tc.priority))
			})
		}

		type popTarget struct {
			priority int64
			value    int
		}
		targets := []popTarget{
			{
				priority: -1,
				value:    500,
			},
			{
				priority: 0,
				value:    300,
			},
			{
				priority: 1,
				value:    100,
			},
			{
				priority: 2,
				value:    200,
			},
			{
				priority: 100,
				value:    400,
			},
		}
		g := NewWithT(t)

		for _, tc := range targets {

			e := pq.Pop()
			g.Expect(e).ToNot(BeNil())
			g.Expect(e.Value).To(Equal(tc.value))
			g.Expect(e.Priority()).To(Equal(tc.priority))
		}
	})
}

func TestPriorityQueueRemove(t *testing.T) {
	t.Run("queue_match_failed", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		e := pq.Add(1, 1)
		g.Expect(e).ToNot(BeNil())
		g.Expect(e.pq).To(Equal(pq))

		e2 := &Element[int]{
			pq: nil,
		}
		err := pq.Remove(e2)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`MismatchQueue`))
	})

	t.Run("out_of_index", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		e := pq.Add(1, 1)
		g.Expect(e).ToNot(BeNil())
		g.Expect(e.pq).To(Equal(pq))

		e.index = -1
		err := pq.Remove(e)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`OutOfIndex`))

		e.index = 2
		err = pq.Remove(e)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`OutOfIndex`))
	})

	t.Run("match_failed", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		e := pq.Add(1, 1)
		g.Expect(e).ToNot(BeNil())
		g.Expect(e.pq).To(Equal(pq))

		e2 := &Element[int]{
			Value:    e.Value,
			priority: e.priority + 1,
			index:    e.index,
			pq:       e.pq,
		}

		err := pq.Remove(e2)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`MismatchPriority`))
	})

	t.Run("ok", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		assertInitPriorityQueue(t, pq)

		for pq.Size() > 0 {
			// random pick one element to remove
			i := rand.Intn(pq.Size())
			e := pq.e[i]
			err := pq.Remove(e)
			g.Expect(err).ToNot(HaveOccurred())

			for _, v := range pq.e {
				g.Expect(v).ToNot(Equal(e))
			}

			assertIsPriorityQueue(t, pq)
		}
	})
}

func assertInitPriorityQueue(t *testing.T, pq *priorityQueue[int]) {
	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority -1",
			value:       500,
			priority:    -1,
			index:       0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			e := pq.Add(tc.value, tc.priority)
			g.Expect(e.Value).To(Equal(tc.value))
			g.Expect(e.index).To(Equal(tc.index))
			g.Expect(e.Priority()).To(Equal(tc.priority))
		})
	}
}

func assertIsPriorityQueue(t *testing.T, pq *priorityQueue[int]) {
	g := NewWithT(t)

	for i, e := range pq.e {
		g.Expect(e.index).To(Equal(i))
	}

	es := make([]*Element[int], len(pq.e))
	copy(es, pq.e)

	heap.Init(pq.h)

	for i := range es {
		g.Expect(es[i]).To(Equal(pq.e[i]))
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	t.Run("queue_match_failed", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		e := pq.Add(1, 1)
		g.Expect(e).ToNot(BeNil())
		g.Expect(e.pq).To(Equal(pq))

		e2 := &Element[int]{
			pq: nil,
		}
		err := pq.Update(e2, e.priority+1)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`MismatchQueue`))
	})

	t.Run("out_of_index", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		e := pq.Add(1, 1)
		g.Expect(e).ToNot(BeNil())
		g.Expect(e.pq).To(Equal(pq))

		e.index = -1
		err := pq.Update(e, e.priority+1)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`OutOfIndex`))

		e.index = 2
		err = pq.Update(e, e.priority+1)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`OutOfIndex`))
	})

	t.Run("ok", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		assertInitPriorityQueue(t, pq)

		c := 100
		for c > 0 {
			c--
			// random pick one element to remove
			i := rand.Intn(pq.Size())
			e := pq.e[i]
			err := pq.Update(e, e.priority+1)
			g.Expect(err).ToNot(HaveOccurred())

			assertIsPriorityQueue(t, pq)
		}
	})

	t.Run("no_change", func(t *testing.T) {
		g := NewWithT(t)
		pq := NewPriorityQueue[int](0).(*priorityQueue[int])

		e := pq.Add(1, 1)
		g.Expect(e).ToNot(BeNil())
		g.Expect(e.pq).To(Equal(pq))

		err := pq.Update(e, e.priority)
		g.Expect(err).ToNot(HaveOccurred())
	})
}

func TestPriorityQueueSize(t *testing.T) {
	pq := NewPriorityQueue[int](0).(*priorityQueue[int])

	type testCase struct {
		description string
		value       int
		priority    int64
		index       int
	}
	testCases := []testCase{
		{
			description: "add priority 1",
			value:       100,
			priority:    1,
			index:       0,
		},
		{
			description: "add priority 2",
			value:       200,
			priority:    2,
			index:       1,
		},
		{
			description: "add priority 0",
			value:       300,
			priority:    0,
			index:       0,
		},
		{
			description: "add priority 100",
			value:       400,
			priority:    100,
			index:       3,
		},
		{
			description: "add priority -1",
			value:       500,
			priority:    -1,
			index:       0,
		},
	}
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			g := NewWithT(t)

			e := pq.Add(tc.value, tc.priority)
			g.Expect(e.Value).To(Equal(tc.value))
			g.Expect(e.index).To(Equal(tc.index))
			g.Expect(e.Priority()).To(Equal(tc.priority))

			g.Expect(pq.Size()).To(Equal(i + 1))
		})
	}
}
