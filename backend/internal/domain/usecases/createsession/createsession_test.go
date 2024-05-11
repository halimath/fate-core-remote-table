package createsession

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

func TestCreateSession(t *testing.T) {
	repo := &repositorymock.Mock{}
	createSession := Provide(repo)

	t.Run("no_user", func(t *testing.T) {
		ctx := context.Background()
		got, err := createSession(ctx, Request{})

		expect.That(t,
			is.Error(err, domain.ErrForbidden),
			is.DeepEqualTo(got, session.Session{}),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		got, err := createSession(ctx, Request{
			Title: "Test",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, session.Session{OwnerID: "2", Title: "Test"}, is.ExcludeFields{"ID"}),
		)
	})
}
