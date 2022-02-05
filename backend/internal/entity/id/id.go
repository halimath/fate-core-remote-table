// Package id defines an ID datatype based on uuid v4
package id

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// ID defines the type for an unique identifier.
type ID string

const (
	urlFriendlyIDAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	urlFriendlyIDLength   = 16
)

var (
	urlFriendlyAlphabetLength = big.NewInt(int64(len(urlFriendlyIDAlphabet)))
)

// New create a new random ID.
func New() ID {
	return ID(uuid.New().String())
}

func NewURLFriendly() ID {
	var buf strings.Builder

	for i := 0; i < urlFriendlyIDLength; i++ {
		idx, err := rand.Int(rand.Reader, urlFriendlyAlphabetLength)
		if err != nil {
			panic(err)
		}
		buf.WriteString(urlFriendlyIDAlphabet[idx.Int64() : idx.Int64()+1])
	}

	return ID(buf.String())
}

// FromString parses s into an ID.
func FromString(s string) ID {
	return ID(s)
}

// String returns a string representation of i.
func (i ID) String() string {
	return string(i)
}
