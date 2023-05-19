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
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIndirectToSetableValue(t *testing.T) {
	type testStruct struct {
		V *int
	}

	type testCase struct {
		desp   string
		v      interface{}
		err    string
		expect interface{}
	}
	testCases := []testCase{
		{
			desp: "normal value",
			v: func() interface{} {
				var i int
				return &i
			}(),
			expect: func() interface{} {
				var i int
				return i
			}(),
		},
		{
			desp: "nil pointer",
			v: func() interface{} {
				st := &testStruct{}
				return reflect.ValueOf(st).Elem().Field(0)
			}(),
			expect: func() interface{} {
				var i int
				return i
			}(),
		},
		{
			desp: "cannot set",
			v:    int(0),
			err:  "cannot been set, it must setable",
		},
		{
			desp: "cannot set pointer",
			v:    (*int)(nil),
			err:  "cannot been set, it must setable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual, err := IndirectToSetableValue(tc.v)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}

func TestIndirectToValue(t *testing.T) {
	type testCase struct {
		desp   string
		v      interface{}
		expect interface{}
	}
	testCases := []testCase{
		{
			desp:   "normal value",
			v:      int(100),
			expect: int(100),
		},
		{
			desp:   "indirect reflect.Value",
			v:      reflect.ValueOf(int(100)),
			expect: int(100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			actual := IndirectToValue(tc.v)
			g.Expect(actual.Interface()).To(Equal(tc.expect))
		})
	}
}
