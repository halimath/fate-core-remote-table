package deleteaspect

import (
	"context"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/usecase"
)

type (
	Request struct {
		SessionID, AspectID string
	}

	Port usecase.Proc[Request]
)

func Provide(r repository.Port) Port {
	return func(ctx context.Context, req Request) error {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return domain.ErrForbidden
		}

		return r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, domain.ErrNotFound
			}

			if s.OwnerID != userID {
				return s, domain.ErrForbidden
			}

			if ok := s.RemoveAspect(req.AspectID); ok {
				return s, nil
			}

			for i := range s.Characters {
				if ok := s.Characters[i].RemoveAspect(req.AspectID); ok {
					return s, nil
				}
			}

			return s, domain.ErrNotFound
		})
	}
}
