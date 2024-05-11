package deleteaspect

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

func TestDeleteAspect(t *testing.T) {
	repo := &repositorymock.Mock{
		S: session.Session{
			ID:      "1",
			OwnerID: "2",
			Characters: []session.Character{
				{
					ID:      "3",
					OwnerID: "4",
					Aspects: session.Aspects{
						{ID: "5"},
					},
				},
			},
			Aspects: session.Aspects{
				{ID: "6"},
			},
		},
	}

	deleteAspect := Provide(repo)

	t.Run("not_authorized", func(t *testing.T) {
		err := deleteAspect(context.Background(), Request{
			SessionID: "1",
			AspectID:  "6",
		})
		expect.That(t, is.Error(err, domain.ErrForbidden))
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, Request{
			SessionID: "2",
			AspectID:  "6",
		})
		expect.That(t, is.Error(err, domain.ErrNotFound))
	})

	t.Run("not_gm", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "4")
		err := deleteAspect(ctx, Request{
			SessionID: "1",
			AspectID:  "6",
		})
		expect.That(t, is.Error(err, domain.ErrForbidden))
	})

	t.Run("aspect_not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, Request{
			SessionID: "1",
			AspectID:  "7",
		})
		expect.That(t, is.Error(err, domain.ErrNotFound))
	})

	t.Run("success_global_aspect", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, Request{
			SessionID: "1",
			AspectID:  "6",
		})
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(repo.S.Aspects, 0),
		)
	})

	t.Run("success_character_aspect", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, Request{
			SessionID: "1",
			AspectID:  "5",
		})
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(repo.S.Characters[0].Aspects, 0),
		)
	})
}
