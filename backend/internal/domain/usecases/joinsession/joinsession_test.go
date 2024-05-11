package joinsession

import (
	"context"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository/repositorymock"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
)

func TestJoinSession(t *testing.T) {
	repo := &repositorymock.Mock{
		S: session.Session{
			ID:      "1",
			OwnerID: "2",
		},
	}
	joinSession := Provide(repo)

	t.Run("not_authorized", func(t *testing.T) {
		characterID, err := joinSession(context.Background(), Request{
			SessionID:     "1",
			CharacterName: "Test",
		})

		expect.That(t,
			is.Error(err, domain.ErrForbidden),
			is.EqualTo(characterID, ""),
			is.SliceOfLen(repo.S.Characters, 0),
		)
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		characterID, err := joinSession(ctx, Request{
			SessionID:     "99",
			CharacterName: "Test",
		})

		expect.That(t,
			is.Error(err, domain.ErrNotFound),
			is.EqualTo(characterID, ""),
			is.SliceOfLen(repo.S.Characters, 0),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		characterID, err := joinSession(ctx, Request{
			SessionID:     "1",
			CharacterName: "Test",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(repo.S, session.Session{
				ID:      "1",
				OwnerID: "2",
				Characters: []session.Character{
					{
						ID:      characterID,
						OwnerID: "3",
						Name:    "Test",
						Type:    session.PC,
					},
				},
			}),
		)
	})
}
