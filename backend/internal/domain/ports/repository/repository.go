// Package repository defines the port interface for a session repository to
// abstract persistence access for a session.
package repository

import (
	"context"
	"errors"

	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
)

// NoSave is a sentinel value returned to signal, that a repository must not be
// saved in the datastore.
var NoSave = errors.New("no save")

// UnitOfWork defines a function type that defines the interface for callback
// functions that perform some atomic unit of work on a [session.Session].
// ctx is used to define deadlines and additional context values. exists signals
// whether the s exists in the datastore. s contains the current values for the
// session.
//
// The function returns a session.Session to update the datastore with if the
// returned error is nil. If the error is the sentinal value NoSave, then no
// insert/update is exectued and the returned session is ignored. Any other
// non-nil error value causes no update and is propagated to any caller.
type UnitOfWork func(ctx context.Context, exists bool, s session.Session) (session.Session, error)

// Port defines the port interface for repository implementations
// managing a session.
type Port interface {
	// Perform executes uow. It first tries to load the session.Session identified
	// by id and passes that information to uow. It returns an error returned
	// from uow or from the underlying datastore.
	Perform(ctx context.Context, id string, uow UnitOfWork) error
}
