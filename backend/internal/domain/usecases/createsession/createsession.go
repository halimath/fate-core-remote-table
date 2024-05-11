package createsession

import (
	"context"
	"fmt"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/usecase"
)

type (
	// Request defines the parameters passed to CreateSession.
	Request struct {
		Title string
	}

	// Port defines the use case type to create a new session.
	Port usecase.Func[Request, session.Session]
)

// Provide provides an implementation for Port use case utilizing r.
func Provide(r repository.Port) Port {
	return func(ctx context.Context, req Request) (ses session.Session, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return session.Session{}, domain.ErrForbidden
		}

		ses = session.Session{
			ID:      id.NewForURL(),
			OwnerID: userID,
			Title:   req.Title,
		}

		err = r.Perform(ctx, ses.ID, func(ctx context.Context, exists bool, _ session.Session) (session.Session, error) {
			if exists {
				return ses, fmt.Errorf("duplicate id: %s", ses.ID)
			}

			return ses, nil
		})
		return
	}
}
