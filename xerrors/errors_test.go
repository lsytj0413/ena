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

package xerrors

import (
	"errors"
	"testing"

	"github.com/onsi/gomega"
)

func TestErrNotFound(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapNotFound("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsNotFound(err)).To(gomega.BeTrue())
	})

	t.Run("is_not_found_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNotFound(nil)).To(gomega.BeFalse())
	})

	t.Run("is_not_found_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		//nolint
		g.Expect(IsNotFound(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrDuplicate(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapDuplicate("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsDuplicate(err)).To(gomega.BeTrue())
	})

	t.Run("is_duplidate_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsDuplicate(nil)).To(gomega.BeFalse())
	})

	t.Run("is_duplidate_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		//nolint
		g.Expect(IsDuplicate(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrNotImplement(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapNotImplement("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsNotImplement(err)).To(gomega.BeTrue())
	})

	t.Run("is_not_implement_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNotImplement(nil)).To(gomega.BeFalse())
	})

	t.Run("is_not_implement_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		//nolint
		g.Expect(IsNotImplement(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrContinue(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapContinue("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsContinue(err)).To(gomega.BeTrue())
	})

	t.Run("is_continue_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsContinue(nil)).To(gomega.BeFalse())
	})

	t.Run("is_continue_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		//nolint
		g.Expect(IsContinue(errors.New("other"))).To(gomega.BeFalse())
	})
}
