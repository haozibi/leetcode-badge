package storage

import (
	"fmt"

	"github.com/pkg/errors"
)

type errNotFound struct {
	e string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("%q not found", e.e)
}

// IsNotFoundError check error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	_, ok := err.(errNotFound)
	return ok
}

// ErrNotFound error not found
func ErrNotFound(s string) error {
	return errNotFound{
		e: s,
	}
}

type errHasExist struct {
	e string
}

func (e errHasExist) Error() string {
	return fmt.Sprintf("%q has exist", e.e)
}

// IsHasExistError check error
func IsHasExistError(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	_, ok := err.(errHasExist)
	return ok
}

// ErrHasExist error: has exist
func ErrHasExist(s string) error {
	return errHasExist{
		e: s,
	}
}
