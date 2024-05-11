// Package domain contains common values used accross the whole domain implementation.
package domain

import "errors"

var (
	// ErrNotFound is a sentinel error value returned when an entity is not found.
	ErrNotFound = errors.New("not found")

	// ErrForbidden is a sentinel error value returned when an operation is forbidden.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidCharacter is a sentinel error value returned when an operation targets a character and that
	// character does not exist or is otherwise invalid.
	ErrInvalidCharacter = errors.New("invlid character")
)
