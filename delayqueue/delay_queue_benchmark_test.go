package delayqueue

import (
	"context"
	"sync"
	"testing"
)

func Benchmark_DelayQueue(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	dq := New[int](0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		dq.Poll(ctx)
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-dq.Chan():
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		dq.Offer(i, defaultTimer.Now()+int64(i))
	}
	cancel()
	wg.Wait()
}

func Benchmark_DelayQueue_Parallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	dq := New[int](0)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		dq.Poll(ctx)
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-dq.Chan():
			}
		}
	}()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			dq.Offer(0, defaultTimer.Now())
		}
	})

	cancel()
	wg.Wait()
}
