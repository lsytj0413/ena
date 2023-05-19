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

	"github.com/lsytj0413/ena"
	"github.com/lsytj0413/ena/xerrors"
)

type getOption struct {
	// Target is the value to store the properties, If the user
	// specified the target, the value will been convert to the target type.
	Target reflect.Value

	// Typ is the target type of value.
	// If the Target field is set, this field will be ignored.
	// Default: string
	Typ reflect.Type

	// Default is the default value (string representation), when the key is not found in the properties,
	// The default is used as value.
	Default *string
}

func defaultGetOption() *getOption {
	return &getOption{
		Target:  reflect.Value{},
		Typ:     reflect.TypeOf(""),
		Default: nil,
	}
}

func (o *getOption) Validate() error {
	if o.Typ == nil {
		return xerrors.Errorf("GetOption.Validate: typ cannot be nil, default to string, did the user override it?")
	}

	if o.IsTargetValid() {
		if o.Target.Type() != o.Typ {
			return xerrors.Errorf(
				"GetOption.Validate: target('%T') doesn't match typ('%v')",
				o.Target.Interface(),
				o.Typ.String(),
			)
		}
	}

	return nil
}

func (o *getOption) Complete() error {
	// First we validate the option, so in the complete workflow, we can do without double check
	if err := o.Validate(); err != nil {
		return err
	}

	// If the target is invalid (user doesn't specify it),
	if !o.IsTargetValid() {
		typ := o.Typ
		if o.Typ.Kind() != reflect.Ptr {
			// If the kind isn't pointer, we should generate pointer for it, to ensure the value is setable
			typ = reflect.PtrTo(typ)
		}

		v := reflect.New(typ.Elem())
		if o.Typ.Kind() != reflect.Ptr {
			v = v.Elem()
		}
		o.Target = v
	}

	// We must validate twice, to ensure after the complete, the option won't violate the rule
	return o.Validate()
}

func (o *getOption) IsTargetValid() bool {
	return o.Target.Kind() != reflect.Invalid
}

// WithDefault will set the default option
func WithDefault(v string) ena.Option[getOption] {
	return ena.NewFnOption(func(opt *getOption) {
		opt.Default = ena.PointerTo(string(v))
	})
}

// WithTarget will set the target option
func WithTarget(i interface{}) ena.Option[getOption] {
	return ena.NewFnOption(func(opt *getOption) {
		opt.Target = ena.IndirectToValue(i)
		opt.Typ = opt.Target.Type()
	})
}

// WithType will set the target type option
func WithType(typ reflect.Type) ena.Option[getOption] {
	return ena.NewFnOption(func(opt *getOption) {
		opt.Typ = typ
	})
}
