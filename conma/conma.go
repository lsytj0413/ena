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
	"sync/atomic"
	"unsafe"

	"github.com/lsytj0413/ena"
	"github.com/lsytj0413/ena/xerrors"
)

// ConfigMgr is the manager for config manager
type ConfigMgr interface {
	// Unmarshal will load the configuration item to object base on
	// tag.
	Unmarshal(o interface{}) error

	// AddConfigReader will add current ConfigReader to list's end.
	AddConfigReader(r ConfigReader)

	// ReadConfig will load all configuration items to current mgr, user
	// must call this before retrieve configuration items.
	ReadConfig() error
}

var defaultConfigMgr unsafe.Pointer

// DefaultConfigMgr will return the default ConfigMgr
func DefaultConfigMgr() ConfigMgr {
	return *(*ConfigMgr)(atomic.LoadPointer(&defaultConfigMgr))
}

// SetDefault will update the default ConfigMgr
func SetDefault(v ConfigMgr) {
	atomic.StorePointer(&defaultConfigMgr, unsafe.Pointer(&v))
}

type configMgr struct {
	configReaders []ConfigReader
	store         *flattenStorage
}

// NewConfigMgr will return an ConfigMgr instance
func NewConfigMgr(readers ...ConfigReader) ConfigMgr {
	return &configMgr{
		store:         newFlattenStorage(),
		configReaders: readers,
	}
}

// Unmarshal will retrieve the configuration items to object o
// base on `conma` struct tag.
func (m *configMgr) Unmarshal(o interface{}) error {
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return xerrors.Errorf("The target object must be an pointer to struct, current is %s", v.String())
	}

	p := &appliers{
		applies: make([]Applier, 0),
	}
	err := m.unmarshalToStruct(v.Elem(), "", p)
	if err != nil {
		return err
	}

	p.Apply()
	return nil
}

func (m *configMgr) unmarshalToStruct(v reflect.Value, prefix string, p *appliers) error {
	if v.Type().Kind() != reflect.Struct {
		panic(xerrors.Errorf("The target value must be an struct, current is %s", v.Type().String()))
	}

	for idx := 0; idx < v.Type().NumField(); idx++ {
		err := m.unmarshalStructField(v, prefix, p, idx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *configMgr) unmarshalStructField(v reflect.Value, prefix string, p *appliers, idx int) error {
	field := v.Type().Field(idx)
	fd, err := NewFieldDescriptor(field, idx)
	if err != nil {
		return err
	}
	if fd == nil {
		return nil
	}

	if fd.Typ.Kind() == reflect.Struct {
		// If this is struct, we dive into sub field
		err := m.unmarshalToStruct(v.Field(idx), fd.Name, p)
		if err != nil {
			return err
		}

		return nil
	}

	applier, err := m.applierForStructField(v, prefix, fd)
	if err != nil {
		return err
	}

	p.Append(applier)
	return nil
}

func (m *configMgr) applierForStructField(v reflect.Value, prefix string, fd *FieldDescriptor) (Applier, error) {
	opts := []ena.Option[getOption]{}
	if fd.Default != nil {
		opts = append(opts, WithDefault(*fd.Default))
	}
	typ := fd.Typ
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PointerTo(typ)
	}
	opts = append(opts, WithType(typ))

	name := fd.Name
	if prefix != "" {
		name = prefix + "." + name
	}

	value, err := m.store.Get(name, opts...)
	if err != nil {
		return nil, err
	}

	vo := reflect.ValueOf(value)
	if fd.Typ.Kind() != reflect.Ptr {
		vo = vo.Elem()
	}

	fv := v.Field(fd.FieldIndex)
	if !fd.Unexported {
		return &fnApplier{
			fn: func() {
				fv.Set(vo)
			},
		}, nil
	}

	return &fnApplier{
		fn: func() {
			// The field is unexported (private field), we cannot directly use Set
			// It will panic with 'reflect.Value.Set using unaddressable value'
			reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem().Set(vo)
		},
	}, nil
}

// AddConfigReader will add current ConfigReader to list's end.
func (m *configMgr) AddConfigReader(r ConfigReader) {
	m.configReaders = append(m.configReaders, r)
}

// ReadConfig will load all configuration items to current mgr, user
// must call this before retrieve configuration items.
func (m *configMgr) ReadConfig() error {
	for _, r := range m.configReaders {
		if err := r.ReadTo(m.store); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	SetDefault(NewConfigMgr())
}
