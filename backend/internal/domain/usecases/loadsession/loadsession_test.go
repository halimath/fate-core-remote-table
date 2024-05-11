package loadsession

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

func TestLoadSession(t *testing.T) {
	repo := &repositorymock.Mock{
		S: session.Session{
			ID:      "1",
			OwnerID: "2",
		},
	}
	loadSession := Provide(repo)

	t.Run("found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		got, err := loadSession(ctx, "1")
		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, session.Session{ID: "1", OwnerID: "2"}),
		)
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		got, err := loadSession(ctx, "2")

		expect.That(t,
			is.Error(err, domain.ErrNotFound),
			is.DeepEqualTo(got, session.Session{}),
		)
	})

	t.Run("no_user", func(t *testing.T) {
		ctx := context.Background()
		got, err := loadSession(ctx, "1")

		expect.That(t,
			is.Error(err, domain.ErrForbidden),
			is.DeepEqualTo(got, session.Session{}),
		)
	})

	t.Run("not_authorized", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")
		got, err := loadSession(ctx, "1")

		expect.That(t,
			is.Error(err, domain.ErrForbidden),
			is.DeepEqualTo(got, session.Session{}),
		)
	})
}
