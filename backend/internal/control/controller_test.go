package control

import (
	"context"
	"testing"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/session"
	"gotest.tools/v3/assert"
)

func TestSessionController_Load(t *testing.T) {
	s := session.New(id.NewURLFriendly(), id.New(), "test")

	c := New()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("found", func(t *testing.T) {
		got, err := c.Load(context.Background(), s.ID)

		assert.NilError(t, err)
		assert.DeepEqual(t, s, got)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := c.Load(context.Background(), id.New())
		assert.DeepEqual(t, session.Session{}, got)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestSessionController_Create(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)

	c := New()

	userID := id.New()
	title := "test"

	sessionID, err := c.Create(context.Background(), userID, title)
	assert.NilError(t, err)

	got, err := c.Load(context.Background(), sessionID)
	assert.NilError(t, err)

	assert.DeepEqual(t, session.Session{
		ID:           sessionID,
		LastModified: now,
		OwnerID:      userID,
		Title:        title,
		Characters:   []session.Character{},
		Aspects:      []session.Aspect{},
	}, got)
}

func TestSessionController_CreateAspect(t *testing.T) {
	s := session.New(id.NewURLFriendly(), id.New(), "test")

	c := New()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("session not found", func(t *testing.T) {
		_, err := c.CreateAspect(context.Background(), s.OwnerID, id.New(), "test")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("not the owner", func(t *testing.T) {
		_, err := c.CreateAspect(context.Background(), id.New(), s.ID, "test")
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("valid", func(t *testing.T) {
		_, err := c.CreateAspect(context.Background(), s.OwnerID, s.ID, "test")
		assert.NilError(t, err)
	})
}

func TestSessionController_DeleteAspect(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		globalAspect := s.AddAspect("test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, id.New(), globalAspect.ID)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("not the owner", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		globalAspect := s.AddAspect("test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), id.New(), s.ID, globalAspect.ID)
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("aspect not found", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, s.ID, id.New())
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("global aspect found", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		globalAspect := s.AddAspect("test")
		character := s.AddCharacter(userID, session.PC, "test")
		character.AddAspect("test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, s.ID, globalAspect.ID)

		assert.NilError(t, err)
		assert.Equal(t, len(c.store[s.ID].s.Aspects), 0)
		assert.Equal(t, len(c.store[s.ID].s.Characters[0].Aspects), 1)
	})

	t.Run("player aspect found", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		s.AddAspect("test")
		character := s.AddCharacter(userID, session.PC, "test")
		playerAspect := character.AddAspect("test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteAspect(context.Background(), userID, s.ID, playerAspect.ID)

		assert.NilError(t, err)
		assert.Equal(t, len(c.store[s.ID].s.Aspects), 1)
		assert.Equal(t, len(c.store[s.ID].s.Characters[0].Aspects), 0)
	})
}

func TestSessionController_CreateCharacter(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		c := New()

		_, err := c.CreateCharacter(context.Background(), id.New(), id.New(), session.NPC, "test")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("only owner can add npc", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		_, err := c.CreateCharacter(context.Background(), id.New(), s.ID, session.NPC, "test")
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("create npc", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		characterID, err := c.CreateCharacter(context.Background(), userID, s.ID, session.NPC, "test")
		assert.NilError(t, err)
		assert.Equal(t, len(c.store[s.ID].s.Characters), 1)
		assert.Equal(t, c.store[s.ID].s.Characters[0].ID, characterID)
	})

	t.Run("create pc", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		characterID, err := c.CreateCharacter(context.Background(), id.New(), s.ID, session.PC, "test")
		assert.NilError(t, err)
		assert.Equal(t, len(c.store[s.ID].s.Characters), 1)
		assert.Equal(t, c.store[s.ID].s.Characters[0].ID, characterID)
	})
}

func TestSessionController_DeleteCharacter(t *testing.T) {
	t.Run("session not found", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		s.AddCharacter(userID, session.PC, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), userID, id.New(), id.New())
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("character not found", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		s.AddCharacter(userID, session.PC, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), userID, s.ID, id.New())
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("neither session nor character owner", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		character := s.AddCharacter(userID, session.PC, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), id.New(), s.ID, character.ID)
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("session owner", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		character := s.AddCharacter(userID, session.PC, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), userID, s.ID, character.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(c.store[s.ID].s.Characters), 0)
	})

	t.Run("character owner", func(t *testing.T) {
		userID := id.New()
		s := session.New(id.NewURLFriendly(), userID, "test")
		character := s.AddCharacter(userID, session.PC, "test")

		c := New()
		c.store[s.ID] = &sessionAndLock{s: s}

		err := c.DeleteCharacter(context.Background(), character.OwnerID, s.ID, character.ID)
		assert.NilError(t, err)
		assert.Equal(t, len(c.store[s.ID].s.Characters), 0)
	})
}

func TestSessionController_CreateCharacterAspect(t *testing.T) {
	userID := id.New()
	s := session.New(id.NewURLFriendly(), userID, "test")
	character := s.AddCharacter(userID, session.PC, "test")

	c := New()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("session not found", func(t *testing.T) {
		_, err := c.CreateCharacterAspect(context.Background(), userID, id.New(), character.ID, "test")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("character not found", func(t *testing.T) {
		_, err := c.CreateCharacterAspect(context.Background(), userID, s.ID, id.New(), "test")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("not the owner", func(t *testing.T) {
		_, err := c.CreateCharacterAspect(context.Background(), id.New(), s.ID, character.ID, "test")
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("add aspect", func(t *testing.T) {
		aspectID, err := c.CreateCharacterAspect(context.Background(), userID, s.ID, character.ID, "test")
		assert.NilError(t, err)
		assert.Equal(t, len(s.Characters[0].Aspects), 1)
		assert.Equal(t, s.Characters[0].Aspects[0], session.Aspect{
			ID:   aspectID,
			Name: "test",
		})
	})
}

func TestSessionController_UpdateFatePoints(t *testing.T) {
	userID := id.New()
	s := session.New(id.NewURLFriendly(), userID, "test")
	character := s.AddCharacter(id.New(), session.PC, "test")

	c := New()
	c.store[s.ID] = &sessionAndLock{s: s}

	t.Run("session not found", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), userID, id.New(), character.ID, 1)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("character not found", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), userID, s.ID, id.New(), 1)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("not the owner", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), id.New(), s.ID, character.ID, 1)
		assert.ErrorIs(t, err, ErrForbidden)
	})

	t.Run("session owner", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), userID, s.ID, character.ID, 1)
		assert.NilError(t, err)
		assert.Equal(t, s.Characters[0].FatePoints, 1)
	})

	t.Run("character owner cannot increase", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), character.OwnerID, s.ID, character.ID, 1)
		assert.ErrorIs(t, err, ErrForbidden)
		assert.Equal(t, s.Characters[0].FatePoints, 1)
	})

	t.Run("cannot decrease below zero", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), character.OwnerID, s.ID, character.ID, -2)
		assert.ErrorIs(t, err, ErrForbidden)
		assert.Equal(t, s.Characters[0].FatePoints, 1)
	})

	t.Run("character owner can decrease", func(t *testing.T) {
		err := c.UpdateFatePoints(context.Background(), character.OwnerID, s.ID, character.ID, -1)
		assert.NilError(t, err)
		assert.Equal(t, s.Characters[0].FatePoints, 0)
	})
}
