package createcharacteraspect

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
		SessionID   string
		Name        string
		CharacterID string
	}

	Port usecase.Func[Request, string]
)

func Provide(r repository.Port) Port {
	return func(ctx context.Context, req Request) (aspectID string, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return "", domain.ErrForbidden
		}

		err = r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, domain.ErrNotFound
			}

			if s.OwnerID != userID {
				return s, domain.ErrForbidden
			}

			c := s.FindCharacter(req.CharacterID)
			if c == nil {
				return s, fmt.Errorf("%w: character not found: %s", domain.ErrInvalidCharacter, req.CharacterID)
			}

			aspectID = c.AddAspect(req.Name).ID
			return s, nil
		})

		return
	}
}
