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
	"strings"

	"github.com/lsytj0413/ena"
)

// FieldDescriptor is the descriptor for conma struct field
// The struct tag must format as:
//  1. `conma:"name:=default"` for config field
type FieldDescriptor struct {
	FieldIndex int
	FieldName  string
	Typ        reflect.Type
	Unexported bool

	Name    string
	Default *string
}

// NewFieldDescriptor ...
func NewFieldDescriptor(field reflect.StructField, idx int) (*FieldDescriptor, error) {
	fd := &FieldDescriptor{
		FieldIndex: idx,
		FieldName:  field.Name,
		Typ:        field.Type,
		Unexported: !field.IsExported(),
	}

	tag, ok := field.Tag.Lookup("conma")
	if !ok {
		return nil, nil
	}

	vv := strings.SplitN(tag, ":=", 2)
	fd.Name = vv[0]
	if len(vv) > 1 {
		fd.Default = ena.PointerTo(vv[1])
	}
	return fd, nil
}
