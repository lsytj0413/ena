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
	"testing"

	. "github.com/onsi/gomega"

	"github.com/lsytj0413/ena/xerrors"
)

type testConfigStore struct {
	setFn func(key string, val interface{}) error
}

func (s *testConfigStore) Set(key string, val interface{}) error {
	return s.setFn(key, val)
}

func TestEnvConfigReaderReadTo(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		r := NewEnvConfigReader()
		count := 0
		vals := map[string]interface{}{}
		s := &testConfigStore{
			setFn: func(key string, val interface{}) error {
				count++
				vals[key] = val
				return nil
			},
		}
		environFn = func() []string {
			return []string{
				"1",
				"1=",
				"k1=v1",
				"k1_k2=v2",
				"k3=v3,v4",
			}
		}

		err := r.ReadTo(s)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(count).To(Equal(4))
		g.Expect(vals).To(Equal(map[string]interface{}{
			"1":     "",
			"k1":    "v1",
			"k1.k2": "v2",
			"k3":    []string{"v3", "v4"},
		}))
	})

	t.Run("single value set failed", func(t *testing.T) {
		g := NewWithT(t)

		r := NewEnvConfigReader()
		count := 0
		vals := map[string]interface{}{}
		s := &testConfigStore{
			setFn: func(key string, val interface{}) error {
				count++
				vals[key] = val

				if key == "k1.k2" {
					return xerrors.ErrContinue
				}

				return nil
			},
		}
		environFn = func() []string {
			return []string{
				"1",
				"1=",
				"k1=v1",
				"k1_k2=v2",
				"k3=v3,v4",
			}
		}

		err := r.ReadTo(s)
		g.Expect(err).To(HaveOccurred())
		g.Expect(count).To(Equal(3))
		g.Expect(vals).To(Equal(map[string]interface{}{
			"1":     "",
			"k1":    "v1",
			"k1.k2": "v2",
		}))
	})

	t.Run("array value set failed", func(t *testing.T) {
		g := NewWithT(t)

		r := NewEnvConfigReader()
		count := 0
		vals := map[string]interface{}{}
		s := &testConfigStore{
			setFn: func(key string, val interface{}) error {
				count++
				vals[key] = val

				if key == "k3" {
					return xerrors.ErrContinue
				}

				return nil
			},
		}
		environFn = func() []string {
			return []string{
				"1",
				"1=",
				"k1=v1",
				"k1_k2=v2",
				"k3=v3,v4",
			}
		}

		err := r.ReadTo(s)
		g.Expect(err).To(HaveOccurred())
		g.Expect(count).To(Equal(4))
		g.Expect(vals).To(Equal(map[string]interface{}{
			"1":     "",
			"k1":    "v1",
			"k1.k2": "v2",
			"k3":    []string{"v3", "v4"},
		}))
	})
}
