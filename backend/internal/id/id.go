// Package id defines an ID datatype based on uuid v4
package id

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// New create a new random ID.
func New() string {
	return uuid.NewString()
}

// NewForURL creates a new ID encoded for URL handling.
func NewForURL() string {
	id := uuid.New()
	return base64.RawURLEncoding.EncodeToString(id[:])
}
