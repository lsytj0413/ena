package ena

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/lsytj0413/ena/xerrors"
)

func TestReceiveChannel(t *testing.T) {
	t.Run("ctx_canceled", func(t *testing.T) {
		g := NewWithT(t)
		ch := make(chan interface{}, 1)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := ReceiveChannel[int](ctx, ch)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`canceled`))
	})

	t.Run("wrong_object_type", func(t *testing.T) {
		g := NewWithT(t)
		ch := make(chan interface{}, 1)
		ch <- "1"

		_, err := ReceiveChannel[int](context.Background(), ch)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`unknown type string, expect int or error`))
	})

	t.Run("receive_error", func(t *testing.T) {
		g := NewWithT(t)
		ch := make(chan interface{}, 1)
		ch <- xerrors.ErrContinue

		_, err := ReceiveChannel[int](context.Background(), ch)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`Continue`))
	})

	t.Run("receive_object", func(t *testing.T) {
		g := NewWithT(t)
		ch := make(chan interface{}, 1)
		ch <- 1

		v, err := ReceiveChannel[int](context.Background(), ch)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(v).To(Equal(1))
	})
}
