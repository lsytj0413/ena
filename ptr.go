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
