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
