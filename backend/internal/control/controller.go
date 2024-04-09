// Package control contains the control layer of the application and provides a SessionController that
// provides business operations for a Session.
package control

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type SessionController interface {
	Load(ctx context.Context, sessionID id.ID) (session.Session, error)
	Create(ctx context.Context, userID id.ID, title string) (id.ID, error)
	CreateAspect(ctx context.Context, userID, sessionID id.ID, name string) (id.ID, error)
	DeleteAspect(ctx context.Context, userID, sessionID, aspectID id.ID) error
	CreateCharacter(ctx context.Context, userID, sessionID id.ID, typ session.CharacterType, name string) (id.ID, error)
	DeleteCharacter(ctx context.Context, userID, sessionID, characterID id.ID) error
	CreateCharacterAspect(ctx context.Context, userID, sessionID, characterID id.ID, name string) (id.ID, error)
	UpdateFatePoints(ctx context.Context, userID, sessionID, characterID id.ID, fatePointsDelta int) error
}

type sessionAndLock struct {
	lock sync.RWMutex
	s    session.Session
}

type unitOfWork func(context.Context, *session.Session) error

type sessionController struct {
	lock  sync.RWMutex
	store map[id.ID]*sessionAndLock
}

func (c *sessionController) withSession(ctx context.Context, s *sessionAndLock, updateLastModified bool, uow unitOfWork) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := uow(ctx, &s.s)

	if err == nil && updateLastModified {
		s.s.LastModified = time.Now().UTC().Truncate(time.Millisecond)
	}

	return err
}

func (c *sessionController) withSessionByID(ctx context.Context, sessionID id.ID, updateLastModified bool, uow unitOfWork) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	s, ok := c.store[sessionID]

	if !ok {
		return fmt.Errorf("%w: session with id %s", ErrNotFound, sessionID)
	}

	return c.withSession(ctx, s, updateLastModified, uow)
}

func (c *sessionController) Load(ctx context.Context, sessionID id.ID) (ses session.Session, err error) {
	err = c.withSessionByID(ctx, sessionID, false, func(ctx context.Context, s *session.Session) error {
		ses = *s
		return nil
	})
	return
}

func (c *sessionController) Create(ctx context.Context, userID id.ID, title string) (id.ID, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var sessionID id.ID
	for {
		sessionID = id.NewURLFriendly()
		if _, ok := c.store[sessionID]; !ok {
			break
		}
	}

	s := session.New(sessionID, userID, title)

	c.store[sessionID] = &sessionAndLock{
		s: s,
	}

	return sessionID, nil
}

func (c *sessionController) CreateAspect(ctx context.Context, userID, sessionID id.ID, name string) (aspectID id.ID, err error) {
	err = c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *session.Session) error {
		if s.OwnerID != userID {
			return fmt.Errorf("%w: only owner may add aspect", ErrForbidden)
		}

		aspectID = s.AddAspect(name).ID

		return nil
	})
	return
}

func (c *sessionController) DeleteAspect(ctx context.Context, userID, sessionID, aspectID id.ID) error {
	return c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *session.Session) error {
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

func (c *sessionController) CreateCharacter(ctx context.Context, userID, sessionID id.ID, typ session.CharacterType, name string) (characterID id.ID, err error) {
	err = c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *session.Session) error {
		if typ == session.NPC && s.OwnerID != userID {
			return fmt.Errorf("%w: only owner may create NPC", ErrForbidden)
		}

		characterID = s.AddCharacter(userID, typ, name).ID

		return nil
	})
	return
}

func (c *sessionController) DeleteCharacter(ctx context.Context, userID, sessionID, characterID id.ID) error {
	return c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *session.Session) error {
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

func (c *sessionController) CreateCharacterAspect(ctx context.Context, userID, sessionID, characterID id.ID, name string) (aspectID id.ID, err error) {
	err = c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *session.Session) error {
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

func (c *sessionController) UpdateFatePoints(ctx context.Context, userID, sessionID, characterID id.ID, fatePointsDelta int) error {
	return c.withSessionByID(ctx, sessionID, true, func(ctx context.Context, s *session.Session) error {
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

func New() *sessionController {
	return &sessionController{
		store: make(map[id.ID]*sessionAndLock),
	}
}

func Provide(cfg config.Config) SessionController {
	c := New()

	if cfg.DevMode {
		generateTestData(c)
	}

	return c
}

func generateTestData(c *sessionController) {
	i := id.FromString("3fa85f64-5717-4562-b3fc-2c963f66afa6")
	owner := id.FromString("00000000-0000-0000-0000-000000000000")

	c.store[i] = &sessionAndLock{
		s: session.Session{
			ID:      i,
			OwnerID: owner,
			Title:   "Test data",
		},
	}
}
