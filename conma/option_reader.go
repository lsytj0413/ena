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
	"strings"
)

// optionConfigReader will read command line arguments to config storage.
type optionConfigReader struct {
}

var (
	// for unittest
	argsFn = func() []string {
		return os.Args[1:]
	}
)

// NewOptionConfigReader will return an option reader
func NewOptionConfigReader() ConfigReader {
	return &optionConfigReader{}
}

// ReadTo will implement ConfigReader.ReadTo method, it will
// 1. parse all command line arguments
// 3. persistent configuration items to ConfigStore
func (r *optionConfigReader) ReadTo(store ConfigStorage) error {
	for _, arg := range argsFn() {
		parts := strings.SplitN(arg, "=", 2)
		switch len(parts) {
		case 1:
			key := strings.TrimLeft(parts[0], "-")
			if len(key) == 0 {
				continue
			}

			if err := store.Set(key, true); err != nil {
				return err
			}
		case 2:
			key := strings.TrimLeft(parts[0], "-")
			if len(key) == 0 {
				continue
			}

			rawValues := parts[1]
			values := strings.Split(rawValues, ",")
			if len(values) > 1 {
				err := store.Set(key, values)
				if err != nil {
					return err
				}
			} else {
				err := store.Set(key, rawValues)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
