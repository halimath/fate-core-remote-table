package session

import (
	"testing"

	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
	"gotest.tools/v3/assert"
)

func TestRemoveByID(t *testing.T) {
	a := []Aspect{
		{
			ID: id.FromString("1"),
		},
		{
			ID: id.FromString("2"),
		},
		{
			ID: id.FromString("3"),
		},
	}

	got := removeByID(&a, id.FromString("0"))
	assert.Equal(t, false, got)
	assert.Equal(t, 3, len(a))

	got = removeByID(&a, id.FromString("2"))
	assert.Equal(t, true, got)
	assert.Equal(t, 2, len(a))

	got = removeByID(&a, id.FromString("3"))
	assert.Equal(t, true, got)
	assert.DeepEqual(t, []Aspect{{ID: id.FromString("1")}}, a)
}

func TestSession_AddAspect(t *testing.T) {
	s := New(id.New(), "test")

	s.AddAspect("test")

	assert.Equal(t, len(s.Aspects), 1)
}

func TestSession_RemoveAspect(t *testing.T) {
	s := New(id.New(), "test")
	aspect := s.AddAspect("test")

	s.RemoveAspect(aspect.ID)

	assert.Equal(t, len(s.Aspects), 0)
}

func TestSession_AddCharacter(t *testing.T) {
	userID := id.New()
	s := New(userID, "test")
	character := s.AddCharacter(userID, PC, "test")

	assert.Equal(t, len(s.Characters), 1)
	assert.DeepEqual(t, s.Characters[0], *character)

}

func TestSession_RemoveCharacter(t *testing.T) {
	userID := id.New()
	s := New(userID, "test")
	character := s.AddCharacter(userID, PC, "test")

	ok := s.RemoveCharacter(character.ID)

	assert.Equal(t, ok, true)
	assert.Equal(t, len(s.Characters), 0)
}

func TestSession_FindCharacter(t *testing.T) {
	userID := id.New()
	s := New(userID, "test")
	c1 := s.AddCharacter(userID, PC, "test")
	c2 := s.AddCharacter(userID, PC, "test2")

	assert.DeepEqual(t, *c1, *s.FindCharacter(c1.ID))
	assert.DeepEqual(t, *c2, *s.FindCharacter(c2.ID))
	assert.Assert(t, s.FindCharacter(id.New()) == nil)
}
