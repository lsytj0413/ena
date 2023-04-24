package ena

import (
	"sync/atomic"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWaitGroupWrapperWrap(t *testing.T) {
	g := NewWithT(t)
	var wg WaitGroupWrapper
	var count uint32

	wg.Wrap(func() {
		atomic.AddUint32(&count, 1)
	}, func() {
		atomic.AddUint32(&count, 2)
	})
	wg.Wrap(func() {
		atomic.AddUint32(&count, 3)
	})

	wg.Wait()
	g.Expect(atomic.LoadUint32(&count)).To(Equal(uint32(6)))
}
