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
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/lsytj0413/ena"
	"github.com/lsytj0413/ena/conv"
	"github.com/lsytj0413/ena/xerrors"
)

// flattenStorage is the storage implement for
// configuration items by kv.
type flattenStorage struct {
	items map[string]string
}

func newFlattenStorage() *flattenStorage {
	return &flattenStorage{
		items: map[string]string{},
	}
}

// nolint
func (s *flattenStorage) Get(key string, opts ...ena.Option[getOption]) (ret interface{}, err error) {
	opt := defaultGetOption()
	for _, o := range opts {
		o.Apply(opt)
	}

	if err = opt.Complete(); err != nil {
		return nil, err
	}

	var targetValue reflect.Value
	targetValue, err = ena.IndirectToSetableValue(opt.Target)
	if err != nil {
		return nil, err
	}

	var propValues []string
	switch {
	case targetValue.Kind() == reflect.Slice:
		vstrs, err := s.doGetSlice(key)
		if err != nil {
			if !xerrors.Is(err, xerrors.ErrNotFound) || opt.Default == nil {
				return nil, err
			}

			if *opt.Default == "" {
				vstrs = []string{}
			} else {
				vstrs = strings.Split(*opt.Default, ",")
			}
		}
		propValues = vstrs
	default:
		vstr, err := s.doGet(key)
		if err != nil {
			if !xerrors.Is(err, xerrors.ErrNotFound) || opt.Default == nil {
				return nil, err
			}

			vstr = *opt.Default
		}
		propValues = []string{vstr}
	}

	cv, err := conv.ConvertTo(context.Background(), targetValue.Type(), propValues)
	if err != nil {
		return nil, err
	}
	targetValue.Set(reflect.ValueOf(cv).Convert(targetValue.Type()))
	return opt.Target.Interface(), nil
}

func (s *flattenStorage) doGet(key string) (string, error) {
	val, ok := s.items[key]
	if ok {
		return val, nil
	}

	return "", xerrors.WrapNotFound("property with key='%v' not found", key)
}

type indexString struct {
	i int64
	v string
}

func (s *flattenStorage) doGetSlice(key string) ([]string, error) {
	val, ok := s.items[key]
	if ok {
		return []string{val}, nil
	}

	ret := []indexString{}
	for k, v := range s.items {
		if !strings.HasPrefix(k, key) {
			continue
		}

		k = k[len(key):]
		if !strings.HasPrefix(k, "[") || !strings.HasSuffix(k, "]") {
			continue
		}

		k = k[1 : len(k)-1]
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}

		ret = append(ret, indexString{
			i: i,
			v: v,
		})
	}

	sort.Slice(ret, func(i int, j int) bool {
		return ret[i].i < ret[j].i
	})

	if len(ret) > 0 {
		rr := make([]string, 0, len(ret))
		for _, v := range ret {
			rr = append(rr, v.v)
		}

		return rr, nil
	}

	return nil, xerrors.WrapNotFound("property slice with key='%v' not found", key)
}

func (s *flattenStorage) Set(key string, val interface{}) error {
	switch v := reflect.ValueOf(val); v.Kind() { //nolint
	case reflect.Map:
		// If the val is a map, we expand the val with keys and set it recursive
		for _, k := range v.MapKeys() {
			kstr, err := conv.ToString(k.Interface())
			if err != nil {
				return xerrors.Wrapf(err, "Cannot convert map's key '%v' to string", k)
			}

			kstr = fmt.Sprintf("%s.%s", key, kstr)
			kvalue := v.MapIndex(k).Interface()
			err = s.Set(kstr, kvalue)
			if err != nil {
				return xerrors.Wrapf(err, "Cannot set val for map's key '%v'", kstr)
			}
		}
	case reflect.Array, reflect.Slice:
		// If the val is a array/slice, we expand the val with index and set it recursive
		for i := 0; i < v.Len(); i++ {
			kstr := fmt.Sprintf("%s[%d]", key, i)
			kvalue := v.Index(i).Interface()
			err := s.Set(kstr, kvalue)
			if err != nil {
				return xerrors.Wrapf(err, "Cannot set val for array/slice index's key '%v'", kstr)
			}
		}
	default:
		value, err := conv.ToString(val)
		if err != nil {
			return xerrors.Wrapf(err, "Cannot convert value to string")
		}
		s.items[key] = value
	}

	return nil
}
