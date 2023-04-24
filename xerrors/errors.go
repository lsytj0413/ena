// Package xerrors contains the utility for errors.
package xerrors

import (
	errors "errors"
	"fmt"
)

// Is alias the errors.Is
var Is = errors.Is

// As alias the errors.As
var As = errors.As

// New alias the errors.New
var New = errors.New

// Wrapf alias the Wrapf
func Wrapf(err error, format string, args ...interface{}) error {
	return Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// Errorf alias the Errorf
var Errorf = fmt.Errorf

var (
	// ErrNotImplement defines the function is not impelement
	ErrNotImplement = errors.New("NotImplement")

	// ErrNotFound defines the object is not found
	ErrNotFound = errors.New("ObjectNotFound")

	// ErrDuplicate defines the object is duplicate
	ErrDuplicate = errors.New("DuplicateObject")

	// ErrContinue defines the err can continue
	ErrContinue = errors.New("Continue")
)

// WrapNotFound return the wraped not found error
func WrapNotFound(format string, args ...interface{}) error {
	return Wrapf(ErrNotFound, format, args...)
}

// IsNotFound return true if the error is object not found
func IsNotFound(err error) bool {
	return Is(err, ErrNotFound)
}

// WrapDuplicate return the wraped duplicate error
func WrapDuplicate(format string, args ...interface{}) error {
	return Wrapf(ErrDuplicate, format, args...)
}

// IsDuplicate return true if the error is duplicate
func IsDuplicate(err error) bool {
	return Is(err, ErrDuplicate)
}

// WrapNotImplement return the not implement error
func WrapNotImplement(format string, args ...interface{}) error {
	return Wrapf(ErrNotImplement, format, args...)
}

// IsNotImplement return true if the error is not implement
func IsNotImplement(err error) bool {
	return Is(err, ErrNotImplement)
}

// WrapContinue return the Continue error
func WrapContinue(format string, args ...interface{}) error {
	return Wrapf(ErrContinue, format, args...)
}

// IsContinue return true if the error is continue
func IsContinue(err error) bool {
	return Is(err, ErrContinue)
}
