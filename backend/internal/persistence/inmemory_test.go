package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
)

func TestInMemory(t *testing.T) {
	repo := NewSessionRepository(config.Config{DevMode: true})

	want := session.Session{
		ID:      "1",
		OwnerID: "2",
		Title:   "Test",
	}

	err := repo.Perform(context.Background(), want.ID, func(_ context.Context, exists bool, s session.Session) (session.Session, error) {
		expect.That(t, is.EqualTo(exists, false))
		return want, nil
	})
	expect.That(t, is.NoError(err))

	var lastModified time.Time

	err = repo.Perform(context.Background(), want.ID, func(_ context.Context, exists bool, s session.Session) (session.Session, error) {
		expect.That(t,
			is.EqualTo(exists, true),
			is.DeepEqualTo(s, want, is.ExcludeFields{"LastModified"}),
			is.EqualTo(s.LastModified.After(want.LastModified), true),
		)
		lastModified = s.LastModified
		return want, repository.NoSave
	})
	expect.That(t, is.NoError(err))

	wantErr := errors.New("kaboom")

	err = repo.Perform(context.Background(), want.ID, func(_ context.Context, exists bool, s session.Session) (session.Session, error) {
		expect.That(t,
			is.EqualTo(exists, true),
			is.EqualTo(s.LastModified, lastModified),
		)
		lastModified = s.LastModified
		return want, wantErr
	})
	expect.That(t, is.Error(err, wantErr))

}
