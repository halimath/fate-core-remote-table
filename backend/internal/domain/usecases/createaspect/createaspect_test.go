package createaspect

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

func TestCreateAspect(t *testing.T) {
	repo := &repositorymock.Mock{
		S: session.Session{
			ID:      "1",
			OwnerID: "2",
		},
	}
	createAspect := Provide(repo)

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		aspectID, err := createAspect(ctx, Request{
			SessionID: "2",
			Name:      "Test",
		})

		expect.That(t,
			is.Error(err, domain.ErrNotFound),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.S.Aspects, 0),
		)
	})

	t.Run("not_authorized", func(t *testing.T) {
		aspectID, err := createAspect(context.Background(), Request{
			SessionID: "2",
			Name:      "Test",
		})

		expect.That(t,
			is.Error(err, domain.ErrForbidden),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.S.Aspects, 0),
		)
	})

	t.Run("not_gm", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		aspectID, err := createAspect(ctx, Request{
			SessionID: "1",
			Name:      "Test",
		})

		expect.That(t,
			is.Error(err, domain.ErrForbidden),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.S.Aspects, 0),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		aspectID, err := createAspect(ctx, Request{
			SessionID: "1",
			Name:      "Test",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(repo.S.Aspects, []session.Aspect{
				{
					ID:   aspectID,
					Name: "Test",
				},
			}),
		)
	})
}
