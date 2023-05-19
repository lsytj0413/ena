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

// envConfigReader will read all available and
// related envs to config storage.
type envConfigReader struct {
}

var (
	// for unittest
	environFn = os.Environ
)

// NewEnvConfigReader return an env reader
func NewEnvConfigReader() ConfigReader {
	return &envConfigReader{}
}

// ReadTo will implement ConfigReader.ReadTo method, it will
// 1. read all envs
// 2. filter by user-defined machinism
// 3. persistent env items to ConfigStore
func (r *envConfigReader) ReadTo(store ConfigStorage) error {
	envStrs := environFn()
	for _, envStr := range envStrs {
		kvArr := strings.SplitN(envStr, "=", 2)
		if len(kvArr) != 2 {
			continue
		}

		key := func(key string) string {
			return strings.ReplaceAll(strings.ToLower(key), "_", ".")
		}(kvArr[0])

		values := strings.Split(kvArr[1], ",")
		if len(values) > 1 {
			err := store.Set(key, values)
			if err != nil {
				return err
			}
		} else {
			err := store.Set(key, kvArr[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
