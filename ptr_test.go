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
	"testing"

	. "github.com/onsi/gomega"
)

func TestPointerTo(t *testing.T) {
	g := NewWithT(t)

	v1 := "123"
	g.Expect(PointerTo("123")).To(Equal(&v1))

	v2 := []int{1, 2}
	g.Expect(PointerTo(v2)).To(Equal(&v2))
}

func TestValueTo(t *testing.T) {
	g := NewWithT(t)

	v1 := "123"
	g.Expect(ValueTo(&v1)).To(Equal(v1))

	v2 := []int{1, 2}
	g.Expect(ValueTo(&v2)).To(Equal(v2))

	var v3 *string
	g.Expect(ValueTo(v3)).To(Equal(""))
}

func TestValueToOrDefault(t *testing.T) {
	g := NewWithT(t)

	v1 := "123"
	g.Expect(ValueToOrDefault(&v1, "456")).To(Equal(v1))

	var v3 *string
	g.Expect(ValueToOrDefault(v3, "456")).To(Equal("456"))
}

func TestEnforcePtr(t *testing.T) {
	type testCase struct {
		desp string
		obj  interface{}

		err string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			obj:  &testCase{},
			err:  ``,
		},
		{
			desp: "not pointer",
			obj:  testCase{},
			err:  `expected pointer, but got ena.testCase type: UnexpectPointer`,
		},
		{
			desp: "invalid kind",
			obj:  nil,
			err:  `expected pointer, but got invalid kind: UnexpectPointer`,
		},
		{
			desp: "nil value",
			obj:  (*testCase)(nil),
			err:  `expected pointer, but got nil: UnexpectPointer`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			_, err := EnforcePtr(tc.obj)
			if tc.err != `` {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}
