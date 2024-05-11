package repositorymock

import (
	"context"
	"errors"

	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
)

type Mock struct {
	S session.Session
}

func (r *Mock) Perform(ctx context.Context, id string, uow repository.UnitOfWork) error {
	s, err := uow(ctx, r.S.ID == id, r.S)
	if errors.Is(err, repository.NoSave) {
		return nil
	}

	if err != nil {
		return err
	}

	r.S = s

	return nil
}
