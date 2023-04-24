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

// PointerTo returns a pointer to the passed value.
func PointerTo[T any](v T) *T {
	return &v
}

// ValueTo returns a value for the passed pointer.
// If the pointer is nil, will return the empty value of T.
func ValueTo[T any](p *T) (v T) {
	if p != nil {
		return *p
	}

	return
}

// ValueToOrDefault returns a value for the passed pointer.
// If the pointer is nil, will return the default value.
func ValueToOrDefault[T any](p *T, d T) (v T) {
	if p != nil {
		return *p
	}

	return d
}

var (
	// ErrUnexpectPointer represent the obj is not an expected pointer type or value
	ErrUnexpectPointer = xerrors.Errorf("UnexpectPointer")
)

// EnforcePtr ensures that obj is a pointer of some sort. Returns a reflect.Value
// of the dereferenced pointer, ensuring that it is settable/addressable.
// Returns an error if this is not possible.
func EnforcePtr(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		if v.Kind() == reflect.Invalid {
			return reflect.Value{}, xerrors.Wrapf(ErrUnexpectPointer, "expected pointer, but got invalid kind")
		}
		return reflect.Value{}, xerrors.Wrapf(ErrUnexpectPointer, "expected pointer, but got %v type", v.Type())
	}
	if v.IsNil() {
		return reflect.Value{}, xerrors.Wrapf(ErrUnexpectPointer, "expected pointer, but got nil")
	}
	return v.Elem(), nil
}
