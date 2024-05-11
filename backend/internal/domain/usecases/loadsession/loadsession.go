// Package loadsession defines the types and functions that implement the get session use case.
package loadsession

import (
	"context"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/usecase"
)

type Port usecase.Func[string, session.Session]

// ProvideLoadSession creates a Func to load a session given its ID.
func Provide(r repository.Port) Port {
	return func(ctx context.Context, sessionID string) (ses session.Session, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return session.Session{}, domain.ErrForbidden
		}

		err = r.Perform(ctx, sessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, domain.ErrNotFound
			}

			if !s.IsMember(userID) {
				return s, domain.ErrForbidden
			}

			ses = s
			return s, repository.NoSave
		})
		return
	}
}
