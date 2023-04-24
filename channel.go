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
