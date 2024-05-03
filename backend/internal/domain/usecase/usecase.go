package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/id"
)

var (
	// ErrNotFound is a sentinel error value returned when an entity is not found.
	ErrNotFound = errors.New("not found")

	// ErrForbidden is a sentinel error value returned when an operation is forbidden.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidCharacter is a sentinel error value returned when an operation targets a character and that
	// character does not exist or is otherwise invalid.
	ErrInvalidCharacter = errors.New("invlid character")
)

// UC is a generic function type that is used to define use case functions that
// provide a return value.
type UC[I, O any] func(context.Context, I) (O, error)

// UCNoRet is a generic function type that is used to define use case functions
// that produce no result (apart from error).
type UCNoRet[I any] func(context.Context, I) error

// -- CreateSession
type (
	// CreateSessionRequest defines the parameters passed to CreateSession.
	CreateSessionRequest struct {
		Title string
	}

	// CreateSession defines the use case type to create a new session.
	CreateSession UC[CreateSessionRequest, session.Session]
)

// ProvideCreateSession provides a CreateSession use case utilizing r.
func ProvideCreateSession(r SessionRepository) CreateSession {
	return func(ctx context.Context, req CreateSessionRequest) (ses session.Session, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return session.Session{}, ErrForbidden
		}

		ses = session.Session{
			ID:      id.NewForURL(),
			OwnerID: userID,
			Title:   req.Title,
		}

		err = r.Perform(ctx, ses.ID, func(ctx context.Context, exists bool, _ session.Session) (session.Session, error) {
			if exists {
				return ses, fmt.Errorf("duplicate id: %s", ses.ID)
			}

			return ses, nil
		})
		return
	}
}

// -- LoadSession

// LoadSession defines the use case function for loading a session.
type LoadSession UC[string, session.Session]

// ProvideLoadSession creates a Func to load a session given its ID.
func ProvideLoadSession(r SessionRepository) LoadSession {
	return func(ctx context.Context, sessionID string) (ses session.Session, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return session.Session{}, ErrForbidden
		}

		err = r.Perform(ctx, sessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, ErrNotFound
			}

			if !s.IsMember(userID) {
				return s, ErrForbidden
			}

			ses = s
			return s, NoSave
		})
		return
	}
}

// -- JoinSession

type (
	JoinSessionRequest struct {
		SessionID     string
		CharacterName string
	}

	JoinSession UC[JoinSessionRequest, string]
)

func ProvideJoinSession(r SessionRepository) JoinSession {
	return func(ctx context.Context, req JoinSessionRequest) (characterID string, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return "", ErrForbidden
		}

		err = r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, ErrNotFound
			}

			characterID = s.AddCharacter(userID, session.PC, req.CharacterName).ID
			return s, nil
		})

		return
	}
}

// -- CreateAspect

type (
	CreateAspectRequest struct {
		SessionID string
		Name      string
	}

	CreateAspect UC[CreateAspectRequest, string]
)

func ProvideCreateAspect(r SessionRepository) CreateAspect {
	return func(ctx context.Context, req CreateAspectRequest) (aspectID string, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return "", ErrForbidden
		}

		err = r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, ErrNotFound
			}

			if s.OwnerID != userID {
				return s, ErrForbidden
			}

			aspectID = s.AddAspect(req.Name).ID
			return s, nil
		})

		return
	}
}

// -- DeleteAspect

type (
	DeleteAspectRequest struct {
		SessionID, AspectID string
	}

	DeleteAspect UCNoRet[DeleteAspectRequest]
)

func ProvideDeleteAspect(r SessionRepository) DeleteAspect {
	return func(ctx context.Context, req DeleteAspectRequest) error {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return ErrForbidden
		}

		return r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, ErrNotFound
			}

			if s.OwnerID != userID {
				return s, ErrForbidden
			}

			if ok := s.RemoveAspect(req.AspectID); ok {
				return s, nil
			}

			for i := range s.Characters {
				if ok := s.Characters[i].RemoveAspect(req.AspectID); ok {
					return s, nil
				}
			}

			return s, ErrNotFound
		})
	}
}

// -- CreateCharacterAspect

type (
	CreateCharacterAspectRequest struct {
		CreateAspectRequest
		CharacterID string
	}

	CreateCharacterAspect UC[CreateCharacterAspectRequest, string]
)

func ProvideCreateCharacterAspect(r SessionRepository) CreateCharacterAspect {
	return func(ctx context.Context, req CreateCharacterAspectRequest) (aspectID string, err error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return "", ErrForbidden
		}

		err = r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, ErrNotFound
			}

			if s.OwnerID != userID {
				return s, ErrForbidden
			}

			c := s.FindCharacter(req.CharacterID)
			if c == nil {
				return s, fmt.Errorf("%w: character not found: %s", ErrInvalidCharacter, req.CharacterID)
			}

			aspectID = c.AddAspect(req.Name).ID
			return s, nil
		})

		return
	}
}

// -- UpdateFatePoints

type (
	UpdateFatePointsRequest struct {
		SessionID, CharacterID string
		Delta                  int
	}

	UpdateFatePoints UCNoRet[UpdateFatePointsRequest]
)

func ProvideUpdateFatePoints(r SessionRepository) UpdateFatePoints {
	return func(ctx context.Context, req UpdateFatePointsRequest) error {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return ErrForbidden
		}

		return r.Perform(ctx, req.SessionID, func(ctx context.Context, exists bool, s session.Session) (session.Session, error) {
			if !exists {
				return s, ErrNotFound
			}

			c := s.FindCharacter(req.CharacterID)
			if c == nil {
				return s, fmt.Errorf("%w: character does not exist: %s", ErrInvalidCharacter, req.CharacterID)
			}

			if s.OwnerID != userID {
				if req.Delta != -1 || c.OwnerID != userID {
					return s, ErrForbidden
				}
			}

			c.FatePoints += req.Delta

			return s, nil
		})
	}
}

// --

var NoSave = errors.New("no save")

type UnitOfWork func(context.Context, bool, session.Session) (session.Session, error)

type SessionRepository interface {
	Perform(ctx context.Context, id string, uow UnitOfWork) error
}
