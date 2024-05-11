package updatefatepoints

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

func TestUpdateFatePoints(t *testing.T) {
	repo := &repositorymock.Mock{
		S: session.Session{
			ID:      "1",
			OwnerID: "2",
			Characters: []session.Character{
				{
					ID:      "3",
					OwnerID: "4",
				},
			},
		},
	}
	updateFatePoints := Provide(repo)

	t.Run("not_authorized", func(t *testing.T) {
		err := updateFatePoints(context.Background(), Request{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, domain.ErrForbidden))
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := updateFatePoints(ctx, Request{
			SessionID:   "99",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, domain.ErrNotFound))
	})

	t.Run("character_not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := updateFatePoints(ctx, Request{
			SessionID:   "1",
			CharacterID: "4",
			Delta:       1,
		})

		expect.That(t, is.Error(err, domain.ErrInvalidCharacter))
	})

	t.Run("neither_gm_nor_player", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "5")
		err := updateFatePoints(ctx, Request{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, domain.ErrForbidden))
	})

	t.Run("player_but_wants_to_increase", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "4")
		err := updateFatePoints(ctx, Request{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, domain.ErrForbidden))
	})

	t.Run("gm_increase", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := updateFatePoints(ctx, Request{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t,
			is.NoError(err),
			is.EqualTo(repo.S.Characters[0].FatePoints, 1),
		)
	})

	t.Run("player_decrease", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "4")
		err := updateFatePoints(ctx, Request{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       -1,
		})

		expect.That(t,
			is.NoError(err),
			is.EqualTo(repo.S.Characters[0].FatePoints, 0),
		)
	})
}
