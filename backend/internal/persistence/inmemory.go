package persistence

import (
	"context"
	"sync"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/domain/ports/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
)

type sessionAndLock struct {
	lock sync.RWMutex
	s    session.Session
}

type repo struct {
	lock  sync.RWMutex
	store map[string]*sessionAndLock
}

func (r *repo) Perform(ctx context.Context, id string, uow repository.UnitOfWork) error {
	r.lock.RLock()
	s, ok := r.store[id]
	r.lock.RUnlock()

	var newSession session.Session
	var err error
	if ok {
		s.lock.Lock()
		defer s.lock.Unlock()
		newSession, err = uow(ctx, ok, s.s)
	} else {
		newSession, err = uow(ctx, ok, session.Session{})
	}

	if err == repository.NoSave {
		return nil
	}

	if err != nil {
		return err
	}

	if s == nil {
		s = &sessionAndLock{}
		r.lock.Lock()
		r.store[newSession.ID] = s
		r.lock.Unlock()
	}

	s.s = newSession
	s.s.LastModified = time.Now().UTC().Truncate(time.Millisecond)

	return nil
}

func NewSessionRepository(cfg config.Config) repository.Port {
	r := &repo{
		store: make(map[string]*sessionAndLock),
	}

	if cfg.DevMode {
		generateTestData(r)
	}

	return r
}

func generateTestData(r *repo) {
	i := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	owner := "00000000-0000-0000-0000-000000000000"

	r.store[i] = &sessionAndLock{
		s: session.Session{
			ID:      i,
			OwnerID: owner,
			Title:   "Test data",
		},
	}
}
