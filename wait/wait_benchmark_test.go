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
