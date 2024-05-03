package id

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func TestNew(t *testing.T) {
	expect.That(t, is.StringOfLen(New(), 36))

}

func TestNewForURL(t *testing.T) {
	expect.That(t, is.StringOfLen(NewForURL(), 22))

}
