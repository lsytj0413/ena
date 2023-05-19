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

package conma

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/lsytj0413/ena"
)

func TestNewFieldDescriptor(t *testing.T) {
	type testCase struct {
		desp  string
		field reflect.StructField
		idx   int

		fd  *FieldDescriptor
		err string
	}
	testCases := []testCase{
		{
			desp: "normal",
			field: reflect.StructField{
				Name: "1",
				Type: reflect.TypeOf("1"),
				Tag:  reflect.StructTag(`conma:"k:=v"`),
			},
			idx: 0,
			fd: &FieldDescriptor{
				FieldIndex: 0,
				FieldName:  "1",
				Typ:        reflect.TypeOf("1"),
				Unexported: false,
				Name:       "k",
				Default:    ena.PointerTo("v"),
			},
			err: ``,
		},
		{
			desp: "didn't have tag",
			field: reflect.StructField{
				Name: "1",
				Type: reflect.TypeOf("1"),
				Tag:  reflect.StructTag(``),
			},
			idx: 0,
			fd:  nil,
			err: ``,
		},
		{
			desp: "didn't have default",
			field: reflect.StructField{
				Name: "1",
				Type: reflect.TypeOf("1"),
				Tag:  reflect.StructTag(`conma:"k"`),
			},
			idx: 0,
			fd: &FieldDescriptor{
				FieldIndex: 0,
				FieldName:  "1",
				Typ:        reflect.TypeOf("1"),
				Unexported: false,
				Name:       "k",
				Default:    nil,
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := NewFieldDescriptor(tc.field, tc.idx)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				g.Expect(actual).To(BeNil())
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.fd))
		})
	}
}
