package joinsession

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
		SessionID     string
		CharacterName string
	}

	Port usecase.Func[Request, string]
)

func Provide(r repository.Port) Port {
	return func(ctx context.Context, req Request) (characterID string, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return "", domain.ErrForbidden
		}

		err = r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, domain.ErrNotFound
			}

			characterID = s.AddCharacter(userID, session.PC, req.CharacterName).ID
			return s, nil
		})

		return
	}
}
