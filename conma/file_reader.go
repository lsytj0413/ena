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

	"gopkg.in/yaml.v3"
)

// fileConfigReader will read all configuration items
// from file to config storage.
type fileConfigReader struct {
	files []string
}

var (
	// for unittest
	readFileFn = os.ReadFile
)

// NewFileConfigReader will return an file reader
func NewFileConfigReader(files ...string) ConfigReader {
	return &fileConfigReader{
		files: files,
	}
}

func (r *fileConfigReader) ReadTo(store ConfigStorage) error {
	for _, file := range r.files {
		data, err := readFileFn(file)
		if err != nil {
			return err
		}

		var m map[string]interface{}
		err = yaml.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		for k, v := range m {
			err = store.Set(k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
