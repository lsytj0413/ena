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

func TestFileConfigReaderReadTo(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		r := NewFileConfigReader("1", "2")
		count := 0
		vals := map[string]interface{}{}
		s := &testConfigStore{
			setFn: func(key string, val interface{}) error {
				count++
				vals[key] = val
				return nil
			},
		}
		readFileFn = func(name string) ([]byte, error) {
			if name == "1" {
				return []byte(`
k1: v1
k2: v2
k3:
  k4: v4`), nil
			}
			return []byte(`
k1: v1
k2: v2.1
k3:
  k4: v4.1
  k5: 
  - v5
  - v5.1`), nil
		}

		err := r.ReadTo(s)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(count).To(Equal(6))
		g.Expect(vals).To(Equal(map[string]interface{}{
			"k1": "v1",
			"k2": "v2.1",
			"k3": map[string]interface{}{
				"k4": "v4.1",
				"k5": []interface{}{
					"v5",
					"v5.1",
				},
			},
		}))
	})

	t.Run("read file failed", func(t *testing.T) {
		g := NewWithT(t)

		r := NewFileConfigReader("1", "2")
		count := 0
		vals := map[string]interface{}{}
		s := &testConfigStore{
			setFn: func(key string, val interface{}) error {
				count++
				vals[key] = val
				return nil
			},
		}
		readFileFn = func(name string) ([]byte, error) {
			if name == "1" {
				return nil, xerrors.ErrContinue
			}
			return []byte(`
k1: v1
k2: v2.1
k3:
  k4: v4.1
  k5: 
  - v5
  - v5.1`), nil
		}

		err := r.ReadTo(s)
		g.Expect(err).To(HaveOccurred())
		g.Expect(count).To(Equal(0))
		g.Expect(vals).To(Equal(map[string]interface{}{}))
	})

	t.Run("unmarshal failed", func(t *testing.T) {
		g := NewWithT(t)

		r := NewFileConfigReader("1", "2")
		count := 0
		vals := map[string]interface{}{}
		s := &testConfigStore{
			setFn: func(key string, val interface{}) error {
				count++
				vals[key] = val
				return nil
			},
		}
		readFileFn = func(name string) ([]byte, error) {
			if name == "1" {
				return []byte(`
k1: v1
k2: v2
k3:
  k4: v4`), nil
			}
			return []byte(`
k1: v1
k2: v2.1
k3:
	k4: v4.1
  k5: 
  - v5
  - v5.1`), nil
		}

		err := r.ReadTo(s)
		g.Expect(err).To(HaveOccurred())
		g.Expect(count).To(Equal(3))
		g.Expect(vals).To(Equal(map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
			"k3": map[string]interface{}{
				"k4": "v4",
			},
		}))
	})
}
