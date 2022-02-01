// Package id defines an ID datatype based on uuid v4
package id

import "github.com/google/uuid"

// ID defines the type for an unique identifier.
type ID string

// New create a new random ID.
func New() ID {
	return ID(uuid.New().String())
}

// FromString parses s into an ID.
func FromString(s string) ID {
	return ID(s)
}

// String returns a string representation of i.
func (i ID) String() string {
	return string(i)
}
