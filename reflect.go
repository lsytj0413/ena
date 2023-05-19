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

	"github.com/lsytj0413/ena/xerrors"
)

// IndirectToValue returns the reflect.Value of v
// If the interface is non reflect.Value, return the reflect.ValueOf(v)
func IndirectToValue(v interface{}) reflect.Value {
	e, ok := v.(reflect.Value)
	if ok {
		return e
	}

	return reflect.ValueOf(v)
}

// IndirectToSetableValue returns the setable reflect.Value of i
func IndirectToSetableValue(i interface{}) (reflect.Value, error) {
	v := IndirectToValue(i)
	if v.Kind() == reflect.Ptr {
		// If the value is pointer and CanSet & IsNil, it maybe the pointer-type field in struct,
		// we must set it with the nil-non-pointer value.
		if v.CanSet() && v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}

	if !v.CanSet() {
		return reflect.Value{}, xerrors.Errorf("The '%T' cannot been set, it must setable", i)
	}
	return v, nil
}
