package session

import (
	"context"
	"testing"
	"time"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fate-core-remote-table/backend/internal/id"
)

func TestService_Load(t *testing.T) {
	s := New(id.NewForURL(), id.New(), "test")

	c := new()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("found", func(t *testing.T) {
		got, err := c.Load(context.Background(), s.ID)

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, s),
		)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := c.Load(context.Background(), id.New())

		expect.That(t,
			is.Error(err, ErrNotFound),
			is.DeepEqualTo(got, Session{}),
		)
	})
}

func TestService_Create(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)

	c := new()

	userID := id.New()
	title := "test"

	sessionID, err := c.Create(context.Background(), userID, title)
	expect.That(t, is.NoError(err))

	got, err := c.Load(context.Background(), sessionID)
	expect.That(t,
		is.NoError(err),
		is.DeepEqualTo(got, Session{
			ID:           sessionID,
			LastModified: now,
			OwnerID:      userID,
			Title:        title,
			Characters:   []Character{},
			Aspects:      []Aspect{},
		}),
	)
}

func TestService_CreateAspect(t *testing.T) {
	s := New(id.NewForURL(), id.New(), "test")

	c := new()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("session not found", func(t *testing.T) {
		_, err := c.CreateAspect(context.Background(), s.OwnerID, id.New(), "test")
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("not the owner", func(t *testing.T) {
		_, err := c.CreateAspect(context.Background(), id.New(), s.ID, "test")
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("valid", func(t *testing.T) {
		_, err := c.CreateAspect(context.Background(), s.OwnerID, s.ID, "test")
		expect.That(t, is.NoError(err))
	})
}

func TestService_DeleteAspect(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		globalAspect := s.AddAspect("test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, id.New(), globalAspect.ID)
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("not the owner", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		globalAspect := s.AddAspect("test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), id.New(), s.ID, globalAspect.ID)
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("aspect not found", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, s.ID, id.New())
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("global aspect found", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		globalAspect := s.AddAspect("test")
		character := s.AddCharacter(userID, PC, "test")
		character.AddAspect("test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, s.ID, globalAspect.ID)

		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(c.store[s.ID].s.Aspects, 0),
			is.SliceOfLen(c.store[s.ID].s.Characters[0].Aspects, 1),
		)
	})

	t.Run("player aspect found", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		s.AddAspect("test")
		character := s.AddCharacter(userID, PC, "test")
		playerAspect := character.AddAspect("test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, s.ID, playerAspect.ID)

		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(c.store[s.ID].s.Aspects, 1),
			is.SliceOfLen(c.store[s.ID].s.Characters[0].Aspects, 0),
		)
	})
}

func TestService_CreateCharacter(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		c := new()

		_, err := c.CreateCharacter(context.Background(), id.New(), id.New(), NPC, "test")
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("only owner can add npc", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		_, err := c.CreateCharacter(context.Background(), id.New(), s.ID, NPC, "test")
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("create npc", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		characterID, err := c.CreateCharacter(context.Background(), userID, s.ID, NPC, "test")
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(c.store[s.ID].s.Characters, 1),
			is.EqualTo(c.store[s.ID].s.Characters[0].ID, characterID),
		)
	})

	t.Run("create pc", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		characterID, err := c.CreateCharacter(context.Background(), id.New(), s.ID, PC, "test")
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(c.store[s.ID].s.Characters, 1),
			is.EqualTo(c.store[s.ID].s.Characters[0].ID, characterID),
		)
	})
}

func TestService_DeleteCharacter(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		s.AddCharacter(userID, PC, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), userID, id.New(), id.New())
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("character not found", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		s.AddCharacter(userID, PC, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), userID, s.ID, id.New())
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("neither session nor character owner", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		character := s.AddCharacter(userID, PC, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), id.New(), s.ID, character.ID)
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("session owner", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		character := s.AddCharacter(userID, PC, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), userID, s.ID, character.ID)
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(c.store[s.ID].s.Characters, 0),
		)
	})

	t.Run("character owner", func(t *testing.T) {
		userID := id.New()
		s := New(id.NewForURL(), userID, "test")
		character := s.AddCharacter(userID, PC, "test")

		c := new()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), character.OwnerID, s.ID, character.ID)
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(c.store[s.ID].s.Characters, 0),
		)
	})
}

func TestService_CreateCharacterAspect(t *testing.T) {
	userID := id.New()
	s := New(id.NewForURL(), userID, "test")
	character := s.AddCharacter(userID, PC, "test")

	c := new()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("session not found", func(t *testing.T) {
		_, err := c.CreateCharacterAspect(context.Background(), userID, id.New(), character.ID, "test")
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("character not found", func(t *testing.T) {
		_, err := c.CreateCharacterAspect(context.Background(), userID, s.ID, id.New(), "test")
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("not the owner", func(t *testing.T) {
		_, err := c.CreateCharacterAspect(context.Background(), id.New(), s.ID, character.ID, "test")
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("add aspect", func(t *testing.T) {
		aspectID, err := c.CreateCharacterAspect(context.Background(), userID, s.ID, character.ID, "test")

		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(s.Characters[0].Aspects, 1),
			is.DeepEqualTo(s.Characters[0].Aspects[0], Aspect{
				ID:   aspectID,
				Name: "test",
			}),
		)
	})
}

func TestService_UpdateFatePoints(t *testing.T) {
	userID := id.New()
	s := New(id.NewForURL(), userID, "test")
	character := s.AddCharacter(id.New(), PC, "test")

	c := new()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("session not found", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), userID, id.New(), character.ID, 1)
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("character not found", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), userID, s.ID, id.New(), 1)
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("not the owner", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), id.New(), s.ID, character.ID, 1)
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("session owner", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), userID, s.ID, character.ID, 1)
		expect.That(t,
			is.NoError(err),
			is.EqualTo(s.Characters[0].FatePoints, 1),
		)
	})

	t.Run("character owner cannot increase", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), character.OwnerID, s.ID, character.ID, 1)
		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(s.Characters[0].FatePoints, 1),
		)
	})

	t.Run("cannot decrease below zero", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), character.OwnerID, s.ID, character.ID, -2)
		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(s.Characters[0].FatePoints, 1),
		)
	})

	t.Run("character owner can decrease", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), character.OwnerID, s.ID, character.ID, -1)
		expect.That(t,
			is.NoError(err),
			is.EqualTo(s.Characters[0].FatePoints, 0),
		)
	})
}
