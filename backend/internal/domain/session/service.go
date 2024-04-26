package session

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type Service interface {
	Load(ctx context.Context, sessionID string) (Session, error)
	Create(ctx context.Context, userID string, title string) (string, error)
	CreateAspect(ctx context.Context, userID, sessionID string, name string) (string, error)
	DeleteAspect(ctx context.Context, userID, sessionID, aspectID string) error
	CreateCharacter(ctx context.Context, userID, sessionID string, typ CharacterType, name string) (string, error)
	DeleteCharacter(ctx context.Context, userID, sessionID, characterID string) error
	CreateCharacterAspect(ctx context.Context, userID, sessionID, characterID string, name string) (string, error)
	UpdateFatePoints(ctx context.Context, userID, sessionID, characterID string, fatePointsDelta int) error
}

type sessionAndLock struct {
	lock sync.RWMutex
	s    Session
}

type unitOfWork func(context.Context, *Session) error

type service struct {
	lock  sync.RWMutex
	store map[string]*sessionAndLock
}

func (c *service) withSession(ctx context.Context, s *sessionAndLock, updateLastModified bool, uow unitOfWork) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := uow(ctx, &s.s)

	if err == nil && updateLastModified {
		s.s.LastModified = time.Now().UTC().Truncate(time.Millisecond)
	}

	return err
}

func (c *service) withSessionByID(ctx context.Context, sessionID string, updateLastModified bool, uow unitOfWork) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	s, ok := c.store[sessionID]

	if !ok {
		return fmt.Errorf("%w: session with id %s", ErrNotFound, sessionID)
	}

	return c.withSession(ctx, s, updateLastModified, uow)
}

func (c *service) Load(ctx context.Context, sessionID string) (ses Session, err error) {
	err = c.withSessionByID(ctx, sessionID, false, func(ctx context.Context, s *Session) error {
		ses = *s
		return nil
	})
	return
}

func (c *service) Create(ctx context.Context, userID string, title string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var sessionID string
	for {
		sessionID = id.NewForURL()
		if _, ok := c.store[sessionID]; !ok {
			break
		}
	}

	s := New(sessionID, userID, title)

	c.store[sessionID] = &sessionAndLock{
		s: s,
	}

	return sessionID, nil
}

func (c *service) CreateAspect(ctx context.Context, userID, sessionID string, name string) (aspectID string, err error) {
	err = c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *Session) error {
		if s.OwnerID != userID {
			return fmt.Errorf("%w: only owner may add aspect", ErrForbidden)
		}

		aspectID = s.AddAspect(name).ID

		return nil
	})
	return
}

func (c *service) DeleteAspect(ctx context.Context, userID, sessionID, aspectID string) error {
	return c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *Session) error {
		if s.OwnerID != userID {
			return fmt.Errorf("%w: only owner may delete aspect", ErrForbidden)
		}

		if s.RemoveAspect(aspectID) {
			return nil
		}

		for i := range s.Characters {
			if s.Characters[i].RemoveAspect(aspectID) {
				return nil
			}
		}

		return fmt.Errorf("%w: aspect with id %s", ErrNotFound, aspectID)
	})
}

func (c *service) CreateCharacter(ctx context.Context, userID, sessionID string, typ CharacterType, name string) (characterID string, err error) {
	err = c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *Session) error {
		if typ == NPC && s.OwnerID != userID {
			return fmt.Errorf("%w: only owner may create NPC", ErrForbidden)
		}

		characterID = s.AddCharacter(userID, typ, name).ID

		return nil
	})
	return
}

func (c *service) DeleteCharacter(ctx context.Context, userID, sessionID, characterID string) error {
	return c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *Session) error {
		for i := range s.Characters {
			if s.Characters[i].ID == characterID {
				if s.OwnerID != userID && s.Characters[i].OwnerID != userID {
					return fmt.Errorf("%w: only owner (session or character) may delete character", ErrForbidden)
				}

				s.RemoveCharacter(characterID)
				return nil
			}
		}

		return fmt.Errorf("%w: character with id %s", ErrNotFound, characterID)
	})
}

func (c *service) CreateCharacterAspect(ctx context.Context, userID, sessionID, characterID string, name string) (aspectID string, err error) {
	err = c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *Session) error {
		if s.OwnerID != userID {
			return fmt.Errorf("%w: only owner may add aspect", ErrForbidden)
		}

		character := s.FindCharacter(characterID)
		if character == nil {
			return fmt.Errorf("%w: character with id %s", ErrNotFound, characterID)
		}

		aspectID = character.AddAspect(name).ID

		return nil
	})
	return
}

func (c *service) UpdateFatePoints(ctx context.Context, userID, sessionID, characterID string, fatePointsDelta int) error {
	return c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *Session) error {
		character := s.FindCharacter(characterID)
		if character == nil {
			return fmt.Errorf("%w: character with id %s", ErrNotFound, characterID)
		}

		if s.OwnerID != userID {
			if character.OwnerID != userID || fatePointsDelta >= 0 {
				return fmt.Errorf("%w: neither session nor character owner", ErrForbidden)
			}
		}

		if character.FatePoints+fatePointsDelta < 0 {
			return fmt.Errorf("%w: cannot reduce fate points below zero", ErrForbidden)
		}

		character.FatePoints += fatePointsDelta

		return nil
	})
}

func new() *service {
	return &service{
		store: make(map[string]*sessionAndLock),
	}
}

func Provide(cfg config.Config) Service {
	srv := new()

	if cfg.DevMode {
		generateTestData(srv)
	}

	return srv
}

func generateTestData(srv *service) {
	i := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	owner := "00000000-0000-0000-0000-000000000000"

	srv.store[i] = &sessionAndLock{
		s: Session{
			ID:      i,
			OwnerID: owner,
			Title:   "Test data",
		},
	}
}
