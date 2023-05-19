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
	"time"

	. "github.com/onsi/gomega"

	"github.com/lsytj0413/ena"
)

func TestFlattenStorageSet(t *testing.T) {
	type testCase struct {
		desp   string
		key    string
		value  interface{}
		err    string
		expect map[string]string
	}
	testCases := []testCase{
		{
			desp:  "normal int8",
			key:   "k1",
			value: int8(1),
			err:   "",
			expect: map[string]string{
				"k1": "1",
			},
		},
		{
			desp:  "normal int16",
			key:   "k1",
			value: int16(0x1000),
			err:   "",
			expect: map[string]string{
				"k1": "4096",
			},
		},
		{
			desp:  "normal int32",
			key:   "k1",
			value: int32(0x10000000),
			err:   "",
			expect: map[string]string{
				"k1": "268435456",
			},
		},
		{
			desp:  "normal int64",
			key:   "k1",
			value: int64(0x1000000000000000),
			err:   "",
			expect: map[string]string{
				"k1": "1152921504606846976",
			},
		},
		{
			desp:  "normal int",
			key:   "k1",
			value: int(0x2000),
			err:   "",
			expect: map[string]string{
				"k1": "8192",
			},
		},
		{
			desp:  "normal uint8",
			key:   "k1",
			value: uint8(1),
			err:   "",
			expect: map[string]string{
				"k1": "1",
			},
		},
		{
			desp:  "normal uint16",
			key:   "k1",
			value: uint16(0x1000),
			err:   "",
			expect: map[string]string{
				"k1": "4096",
			},
		},
		{
			desp:  "normal uint32",
			key:   "k1",
			value: uint32(0x10000000),
			err:   "",
			expect: map[string]string{
				"k1": "268435456",
			},
		},
		{
			desp:  "normal uint64",
			key:   "k1",
			value: uint64(0x1000000000000000),
			err:   "",
			expect: map[string]string{
				"k1": "1152921504606846976",
			},
		},
		{
			desp:  "normal uint",
			key:   "k1",
			value: uint(0x2000),
			err:   "",
			expect: map[string]string{
				"k1": "8192",
			},
		},
		{
			desp:  "normal float32",
			key:   "k1",
			value: float32(1.1),
			err:   "",
			expect: map[string]string{
				"k1": "1.1",
			},
		},
		{
			desp:  "normal float64",
			key:   "k1",
			value: float64(1.2),
			err:   "",
			expect: map[string]string{
				"k1": "1.2",
			},
		},
		{
			desp:  "normal bool",
			key:   "k1",
			value: bool(true),
			err:   "",
			expect: map[string]string{
				"k1": "true",
			},
		},
		{
			desp:  "normal bool",
			key:   "k1",
			value: time.Duration(1) * time.Second,
			err:   "",
			expect: map[string]string{
				"k1": "1s",
			},
		},
		{
			desp:  "normal array",
			key:   "k1",
			value: [3]int{1, 2, 3},
			err:   "",
			expect: map[string]string{
				"k1[0]": "1",
				"k1[1]": "2",
				"k1[2]": "3",
			},
		},
		{
			desp:  "normal slice",
			key:   "k1",
			value: []int{1, 2, 3},
			err:   "",
			expect: map[string]string{
				"k1[0]": "1",
				"k1[1]": "2",
				"k1[2]": "3",
			},
		},
		{
			desp: "normal map",
			key:  "k1",
			value: map[string]int{
				"k1": 1,
				"k2": 2,
			},
			err: "",
			expect: map[string]string{
				"k1.k1": "1",
				"k1.k2": "2",
			},
		},
		{
			desp: "normal map with int key",
			key:  "k1",
			value: map[int]int{
				1: 1,
				2: 2,
			},
			err: "",
			expect: map[string]string{
				"k1.1": "1",
				"k1.2": "2",
			},
		},
		{
			desp:  "set value for struct failed",
			key:   "k1",
			value: testCase{},
			err:   "Cannot convert value to string",
		},
		{
			desp:  "set value for array struct failed",
			key:   "k1",
			value: [1]testCase{{}},
			err:   `Cannot set val for array/slice index's key 'k1\[0\]'`,
		},
		{
			desp:  "set value for slice struct failed",
			key:   "k1",
			value: []testCase{{}},
			err:   `Cannot set val for array/slice index's key 'k1\[0\]'`,
		},
		{
			desp: "set value for map struct failed",
			key:  "k1",
			value: map[string]testCase{
				"k1": {},
			},
			err: "Cannot set val for map's key 'k1.k1'",
		},
		{
			desp: "set value for map key struct failed",
			key:  "k1",
			value: map[*testCase]string{
				{}: "",
			},
			err: "Cannot convert map's key",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := newFlattenStorage()
			err := p.Set(tc.key, tc.value)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).Should(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(map[string]string(p.items)).To(Equal(tc.expect))
		})
	}
}

func TestFlattenStorageGet(t *testing.T) {
	type testStruct struct {
		V string
	}

	type testCase struct {
		desp   string
		p      *flattenStorage
		key    string
		target interface{}
		typ    reflect.Type
		def    *string

		err    string
		expect interface{}
	}
	testCases := []testCase{
		{
			desp: "normal get with non option",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "v1",
				},
			},
			key:    "k1",
			err:    "",
			expect: "v1",
		},
		{
			desp: "key not found without default value",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "v1",
				},
			},
			key:    "k0",
			err:    "property with key='k0' not found",
			expect: "",
		},
		{
			desp: "not found with default value",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "v1",
				},
			},
			key:    "k0",
			def:    ena.PointerTo("1"),
			err:    "",
			expect: "1",
		},
		{
			desp: "already exists with default value",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "v1",
				},
			},
			key:    "k1",
			def:    ena.PointerTo("1"),
			err:    "",
			expect: "v1",
		},
		{
			desp: "normal get uint8",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((uint8)(0)),
			err: "",
			expect: func() interface{} {
				var i uint8 = 1
				return i
			}(),
		},
		{
			desp: "normal get uint16",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "4096",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((uint16)(0)),
			err: "",
			expect: func() interface{} {
				var i uint16 = 4096
				return i
			}(),
		},
		{
			desp: "normal get uint32",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "268435456",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((uint32)(0)),
			err: "",
			expect: func() interface{} {
				var i uint32 = 268435456
				return i
			}(),
		},
		{
			desp: "normal get uint64",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1152921504606846976",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((uint64)(0)),
			err: "",
			expect: func() interface{} {
				var i uint64 = 1152921504606846976
				return i
			}(),
		},
		{
			desp: "normal get uint",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1152921504606846976",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((uint)(0)),
			err: "",
			expect: func() interface{} {
				var i uint = 1152921504606846976
				return i
			}(),
		},
		{
			desp: "normal get int8",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((int8)(0)),
			err: "",
			expect: func() interface{} {
				var i int8 = 1
				return i
			}(),
		},
		{
			desp: "normal get int16",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "4096",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((int16)(0)),
			err: "",
			expect: func() interface{} {
				var i int16 = 4096
				return i
			}(),
		},
		{
			desp: "normal get int32",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "268435456",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((int32)(0)),
			err: "",
			expect: func() interface{} {
				var i int32 = 268435456
				return i
			}(),
		},
		{
			desp: "normal get int64",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1152921504606846976",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((int64)(0)),
			err: "",
			expect: func() interface{} {
				var i int64 = 1152921504606846976
				return i
			}(),
		},
		{
			desp: "normal get int",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1152921504606846976",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((int)(0)),
			err: "",
			expect: func() interface{} {
				var i = 1152921504606846976
				return i
			}(),
		},
		{
			desp: "normal get float32",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1.1",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((float32)(0)),
			err: "",
			expect: func() interface{} {
				var i float32 = 1.1
				return i
			}(),
		},
		{
			desp: "normal get float64",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1.1",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((float64)(0)),
			err: "",
			expect: func() interface{} {
				var i = 1.1
				return i
			}(),
		},
		{
			desp: "normal get bool",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "true",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((bool)(false)),
			err: "",
			expect: func() interface{} {
				var i = true
				return i
			}(),
		},
		{
			desp: "normal get string",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "true",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf((string)("")),
			err: "",
			expect: func() interface{} {
				var i = "true"
				return i
			}(),
		},
		{
			desp: "normal get time duration",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "1s",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf(time.Duration(0)),
			err: "",
			expect: func() interface{} {
				i := time.Duration(1) * time.Second
				return i
			}(),
		},
		{
			desp: "unsupport target type",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "true",
					"k2": "2",
				},
			},
			key: "k1",
			typ: reflect.TypeOf(testCase{}),
			err: "Unsupport target type conma.testCase",
		},
		{
			desp: "target to struct ptr value",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "true",
					"k2": "2",
				},
			},
			key: "k1",
			target: func() interface{} {
				s := &testStruct{}
				return &(s.V)
			}(),
			err: "",
			expect: func() interface{} {
				v := "true"
				return &v
			}(),
		},
		{
			desp: "target to struct field value ptr",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "true1",
					"k2": "2",
				},
			},
			key: "k1",
			target: func() interface{} {
				s := testStruct{}
				return &(s.V)
			}(),
			err: "",
			expect: func() interface{} {
				v := "true1"
				return &v
			}(),
		},
		{
			desp: "all option set",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "false",
					"k2": "2",
				},
			},
			key: "k3",
			target: func() interface{} {
				s := testStruct{}
				return &(s.V)
			}(),
			typ: reflect.TypeOf((*string)(nil)),
			def: func() *string {
				v := "true2"
				return &v
			}(),
			err: "",
			expect: func() interface{} {
				v := "true2"
				return &v
			}(),
		},
		{
			desp: "option validate failed",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "true3",
					"k2": "2",
				},
			},
			key: "k3",
			target: func() interface{} {
				s := testStruct{}
				return &(s.V)
			}(),
			typ: reflect.TypeOf((*int)(nil)),
			def: func() *string {
				v := "true3"
				return &v
			}(),
			err:    `target\('\*string'\) doesn't match typ\('\*int'\)`,
			expect: nil,
		},
		{
			desp: "get []string without default and one key hit",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "false",
					"k2": "2",
				},
			},
			key: "k1",
			target: func() interface{} {
				s := []string{}
				return &s
			}(),
			def: nil,
			err: "",
			expect: func() interface{} {
				v := []string{"false"}
				return &v
			}(),
		},
		{
			desp: "get []string without default and multi key hit",
			p: &flattenStorage{
				items: map[string]string{
					"k1[1]": "false",
					"k1[2]": "true",
					"k2":    "2",
				},
			},
			key: "k1",
			target: func() interface{} {
				s := []string{}
				return &s
			}(),
			def: nil,
			err: "",
			expect: func() interface{} {
				v := []string{"false", "true"}
				return &v
			}(),
		},
		{
			desp: "get []string with default",
			p: &flattenStorage{
				items: map[string]string{
					"k2": "2",
				},
			},
			key: "k1",
			target: func() interface{} {
				s := []string{}
				return &s
			}(),
			def: ena.PointerTo("false,true"),
			err: "",
			expect: func() interface{} {
				v := []string{"false", "true"}
				return &v
			}(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			opts := make([]ena.Option[getOption], 0)
			if tc.target != nil {
				opts = append(opts, WithTarget(tc.target))
			}
			if tc.typ != nil {
				opts = append(opts, WithType(tc.typ))
			}
			if tc.def != nil {
				opts = append(opts, WithDefault(*tc.def))
			}

			actual, err := tc.p.Get(tc.key, opts...)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))

			if tc.target != nil {
				g.Expect(tc.target).To(Equal(actual))
				g.Expect(reflect.TypeOf(tc.target)).To(Equal(reflect.TypeOf(tc.expect)))
			}
		})
	}
}

func TestFlattenStorageDoGetSlice(t *testing.T) {
	type testCase struct {
		desp string
		p    *flattenStorage
		key  string

		err    string
		expect interface{}
	}
	testCases := []testCase{
		{
			desp: "normal exactly match",
			p: &flattenStorage{
				items: map[string]string{
					"k1": "v1",
				},
			},
			key:    "k1",
			err:    "",
			expect: []string{"v1"},
		},
		{
			desp: "normal slice match",
			p: &flattenStorage{
				items: map[string]string{
					"k1[0]": "v1",
					"k1[1]": "v2",
				},
			},
			key:    "k1",
			err:    "",
			expect: []string{"v1", "v2"},
		},
		{
			desp: "index parse failed",
			p: &flattenStorage{
				items: map[string]string{
					"k1[a]": "v1",
					"k1[1]": "v2",
				},
			},
			key:    "k1",
			err:    `parsing "a": invalid syntax`,
			expect: []string{},
		},
		{
			desp: "not found",
			p: &flattenStorage{
				items: map[string]string{
					"k1[0]": "v1",
					"k1[1]": "v2",
				},
			},
			key:    "k2",
			err:    "property slice with key='k2' not found",
			expect: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := tc.p.doGetSlice(tc.key)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
