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
	"context"

	"github.com/lsytj0413/ena/xerrors"
)

// ReceiveChannel consume obj from channel, it will:
// 1. return err if ctx is Done
// 2. return err is obj is error
// 3. return err if obj is not error or type T
func ReceiveChannel[T any](ctx context.Context, ch <-chan interface{}) (v T, err error) {
	select {
	case <-ctx.Done():
		return v, ctx.Err()
	case vv := <-ch:
		switch vvo := vv.(type) {
		case error:
			return v, vvo
		case T:
			return vvo, nil
		}

		return v, xerrors.Errorf("unknown type %T, expect %T or error", vv, v)
	}
}
