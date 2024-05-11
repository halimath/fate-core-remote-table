package updatefatepoints

import (
	"context"
	"fmt"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/usecase"
)

type (
	Request struct {
		SessionID, CharacterID string
		Delta                  int
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

			c := s.FindCharacter(req.CharacterID)
			if c == nil {
				return s, fmt.Errorf("%w: character does not exist: %s", domain.ErrInvalidCharacter, req.CharacterID)
			}

			if s.OwnerID != userID {
				if req.Delta != -1 || c.OwnerID != userID {
					return s, domain.ErrForbidden
				}
			}

			c.FatePoints += req.Delta

			return s, nil
		})
	}
}
