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
	"os"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/lsytj0413/ena"
	"github.com/lsytj0413/ena/xerrors"
)

func TestConfigMgrUnmarshal(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		type testObj struct {
			f1 string
			f2 string        `conma:"f2"`
			F3 int           `conma:"f3"`
			F4 time.Duration `conma:"f4.4:=1m"`
			e1 string        `conma:"e1"`
			o1 string        `conma:"o1"`
			s1 struct {
				f1 string `conma:"f2"`
			} `conma:"f2"`
		}
		g := NewWithT(t)

		o := testObj{}
		os.Setenv("e1", "e1v")
		os.Args = []string{"", "-o1=o1v"}
		f := NewConfigMgr(&envConfigReader{}, &optionConfigReader{}).(*configMgr)
		err := f.ReadConfig()
		g.Expect(err).ToNot(HaveOccurred())

		f.store.Set("f2", "1")
		f.store.Set("f2.f2", "2")
		f.store.Set("f3", "10")
		err = f.Unmarshal(&o)
		g.Expect(err).ToNot(HaveOccurred())

		expect := testObj{
			f1: "",
			f2: "1",
			F3: 10,
			F4: time.Minute,
			e1: "e1v",
			o1: "o1v",
		}
		expect.s1.f1 = "2"
		g.Expect(o).To(Equal(expect))
	})

	t.Run("non pointer to struct", func(t *testing.T) {
		//nolint
		type testObj struct {
			f1 string
			f2 string        `conma:"f2"`
			F3 int           `conma:"f3"`
			F4 time.Duration `conma:"f4.4:=1m"`
			e1 string        `conma:"e1"`
			o1 string        `conma:"o1"`
			s1 struct {
				f1 string `conma:"f2"`
			} `conma:"f2"`
		}
		g := NewWithT(t)

		o := testObj{}
		os.Setenv("e1", "e1v")
		os.Args = []string{"", "-o1=o1v"}
		f := NewConfigMgr(&envConfigReader{}, &optionConfigReader{}).(*configMgr)
		err := f.ReadConfig()
		g.Expect(err).ToNot(HaveOccurred())

		f.store.Set("f2", "1")
		f.store.Set("f2.f2", "2")
		f.store.Set("f3", "10")
		err = f.Unmarshal(o.e1)
		g.Expect(err).To(HaveOccurred())
	})

	t.Run("failed test", func(t *testing.T) {
		//nolint
		type testObj struct {
			f1 string
			f2 string        `conma:"f2"`
			F3 int           `conma:"f3"`
			F4 time.Duration `conma:"f4.4:=1m"`
			e1 string        `conma:"e1"`
			o1 string        `conma:"o1"`
			s1 struct {
				f1 string `conma:"f2"`
			} `conma:"f2"`
		}
		g := NewWithT(t)

		o := testObj{}
		os.Setenv("e1", "e1v")
		os.Args = []string{"", "-o1=o1v"}
		f := NewConfigMgr(&envConfigReader{}, &optionConfigReader{}).(*configMgr)
		err := f.ReadConfig()
		g.Expect(err).ToNot(HaveOccurred())

		f.store.Set("f2", "1")
		f.store.Set("f3", "10")
		err = f.Unmarshal(&o)
		g.Expect(err).To(HaveOccurred())
	})
}

type testConfigReader struct {
	fnReadTo func(r ConfigStorage) error
}

func (r *testConfigReader) ReadTo(store ConfigStorage) error {
	return r.fnReadTo(store)
}

func TestConfigMgrReadConfig(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		count := 0
		mgr := NewConfigMgr(
			&testConfigReader{
				fnReadTo: func(r ConfigStorage) error {
					count++
					return nil
				},
			},
			&testConfigReader{
				fnReadTo: func(r ConfigStorage) error {
					count++
					return nil
				},
			},
			&testConfigReader{
				fnReadTo: func(r ConfigStorage) error {
					count++
					return nil
				},
			},
		)
		err := mgr.ReadConfig()
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(count).To(Equal(3))
	})

	t.Run("fail test", func(t *testing.T) {
		g := NewWithT(t)

		count := 0
		mgr := NewConfigMgr(
			&testConfigReader{
				fnReadTo: func(r ConfigStorage) error {
					count++
					return nil
				},
			},
			&testConfigReader{
				fnReadTo: func(r ConfigStorage) error {
					count++
					return xerrors.ErrContinue
				},
			},
			&testConfigReader{
				fnReadTo: func(r ConfigStorage) error {
					count++
					return nil
				},
			},
		)
		err := mgr.ReadConfig()
		g.Expect(err).To(HaveOccurred())
		g.Expect(count).To(Equal(2))
	})
}

func TestConfigMgrAddConfigReader(t *testing.T) {
	g := NewWithT(t)
	mgr := NewConfigMgr()

	mgr.AddConfigReader(&testConfigReader{})
	g.Expect(mgr.(*configMgr).configReaders).To(HaveLen(1))
}

func TestConfigMgrUnmarshalToStruct(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"p.k1": "v1",
				"k2":   "v2",
			},
		}

		type S1 struct {
			v1 string `conma:"k1"`
		}
		type S2 struct {
			v2 string `conma:"k2"`
			s1 S1     `conma:"p"`
			v3 string
		}

		var target S2
		p := &appliers{}
		err := mgr.unmarshalToStruct(reflect.ValueOf(&target).Elem(), "", p)
		g.Expect(err).ToNot(HaveOccurred())
		p.Apply()
		g.Expect(target.s1.v1).To(Equal("v1"))
		g.Expect(target.v2).To(Equal("v2"))
		g.Expect(target.v3).To(Equal(""))
	})

	t.Run("failed test", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"p.k1": "v1",
			},
		}

		type S1 struct {
			v1 string `conma:"k1"`
		}
		type S2 struct {
			s1 S1     `conma:"p"`
			v2 string `conma:"k2"`
		}

		var target S2
		p := &appliers{}
		err := mgr.unmarshalToStruct(reflect.ValueOf(&target).Elem(), "", p)
		g.Expect(err).To(HaveOccurred())
		p.Apply()
		g.Expect(target.s1.v1).To(Equal("v1"))
		g.Expect(target.v2).To(Equal(""))
	})

	t.Run("non-struct", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"p.k1": "v1",
			},
		}

		var target string
		p := &appliers{}
		g.Ω(func() {
			mgr.unmarshalToStruct(reflect.ValueOf(&target).Elem(), "", p)
		}).Should(Panic())
	})
}

func TestConfigMgrUnmarshalStructField(t *testing.T) {
	t.Run("normal test with struct field", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"p.k1": "v1",
				"k2":   "v2",
			},
		}

		type S1 struct {
			v1 string `conma:"k1"`
		}
		type S2 struct {
			s1 S1 `conma:"p"`
		}

		var target S2
		p := &appliers{}
		err := mgr.unmarshalStructField(reflect.ValueOf(&target).Elem(), "", p, 0)
		g.Expect(err).ToNot(HaveOccurred())
		p.Apply()
		g.Expect(target.s1.v1).To(Equal("v1"))
	})

	t.Run("normal test with string field", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"p.k1": "v1",
				"k2":   "v2",
			},
		}

		type S2 struct {
			v2 string `conma:"k2"`
		}

		var target S2
		p := &appliers{}
		err := mgr.unmarshalStructField(reflect.ValueOf(&target).Elem(), "", p, 0)
		g.Expect(err).ToNot(HaveOccurred())
		p.Apply()
		g.Expect(target.v2).To(Equal("v2"))
	})

	t.Run("failed test with struct field", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"k2": "v2",
			},
		}

		type S1 struct {
			v1 string `conma:"k1"`
		}
		type S2 struct {
			s1 S1 `conma:"p"`
		}

		var target S2
		p := &appliers{}
		err := mgr.unmarshalStructField(reflect.ValueOf(&target).Elem(), "", p, 0)
		g.Expect(err).To(HaveOccurred())
		p.Apply()
		g.Expect(target.s1.v1).To(Equal(""))
	})

	t.Run("failed test with string field", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"p.k1": "v1",
			},
		}

		type S2 struct {
			v2 string `conma:"k2"`
		}

		var target S2
		p := &appliers{}
		err := mgr.unmarshalStructField(reflect.ValueOf(&target).Elem(), "", p, 0)
		g.Expect(err).To(HaveOccurred())
		p.Apply()
		g.Expect(target.v2).To(Equal(""))
	})
}

func TestConfigMgrApplierForStructField(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"k1.k2": "v1",
			},
		}

		type targetS struct {
			v string
		}
		var target targetS
		applier, err := mgr.applierForStructField(reflect.ValueOf(&target).Elem(), "k1", &FieldDescriptor{
			Default:    nil,
			Typ:        reflect.TypeOf(""),
			FieldIndex: 0,
			Unexported: true,
			Name:       "k2",
		})
		g.Expect(err).ToNot(HaveOccurred())

		applier.Apply()
		g.Expect(target.v).To(Equal("v1"))
	})

	t.Run("normal wit default & exported", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"k1.k2": "v1",
			},
		}

		type targetS struct {
			V string
		}
		var target targetS
		applier, err := mgr.applierForStructField(reflect.ValueOf(&target).Elem(), "k1", &FieldDescriptor{
			Default:    ena.PointerTo("v3"),
			Typ:        reflect.TypeOf(""),
			FieldIndex: 0,
			Unexported: false,
			Name:       "k3",
		})
		g.Expect(err).ToNot(HaveOccurred())

		applier.Apply()
		g.Expect(target.V).To(Equal("v3"))
	})

	t.Run("normal test with pointer", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"k1.k2": "v1",
			},
		}

		type targetS struct {
			v *string
		}
		var target targetS
		applier, err := mgr.applierForStructField(reflect.ValueOf(&target).Elem(), "k1", &FieldDescriptor{
			Default:    nil,
			Typ:        reflect.PtrTo(reflect.TypeOf("")),
			FieldIndex: 0,
			Unexported: true,
			Name:       "k2",
		})
		g.Expect(err).ToNot(HaveOccurred())

		applier.Apply()
		g.Expect(target.v).To(Equal(ena.PointerTo("v1")))
	})

	t.Run("get value from store failed", func(t *testing.T) {
		g := NewWithT(t)
		mgr := NewConfigMgr().(*configMgr)
		mgr.store = &flattenStorage{
			items: map[string]string{
				"k1.k3": "v1",
			},
		}

		type targetS struct {
			v string //nolint
		}
		var target targetS
		applier, err := mgr.applierForStructField(reflect.ValueOf(&target).Elem(), "k1", &FieldDescriptor{
			Default:    nil,
			Typ:        reflect.TypeOf(""),
			FieldIndex: 0,
			Unexported: true,
			Name:       "k2",
		})
		g.Expect(err).To(HaveOccurred())
		g.Expect(applier).To(BeNil())
	})
}
