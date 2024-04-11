package session

import (
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
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
	expect.That(t,
		is.EqualTo(got, false),
		is.SliceOfLen(a, 3),
	)

	got = removeByID(&a, id.FromString("2"))
	expect.That(t,
		is.EqualTo(got, true),
		is.SliceOfLen(a, 2),
	)

	got = removeByID(&a, id.FromString("3"))
	expect.That(t,
		is.EqualTo(got, true),
		is.DeepEqualTo(a, []Aspect{{ID: id.FromString("1")}}),
	)
}

func TestSession_AddAspect(t *testing.T) {
	s := New(id.NewURLFriendly(), id.New(), "test")

	s.AddAspect("test")

	expect.That(t,
		is.SliceOfLen(s.Aspects, 1),
	)
}

func TestSession_RemoveAspect(t *testing.T) {
	s := New(id.NewURLFriendly(), id.New(), "test")
	aspect := s.AddAspect("test")

	s.RemoveAspect(aspect.ID)

	expect.That(t,
		is.SliceOfLen(s.Aspects, 0),
	)
}

func TestSession_AddCharacter(t *testing.T) {
	userID := id.New()
	s := New(id.NewURLFriendly(), userID, "test")
	character := s.AddCharacter(userID, PC, "test")

	expect.That(t,
		is.DeepEqualTo(s.Characters, []Character{
			{
				ID:      character.ID,
				OwnerID: userID,
				Name:    "test",
				Type:    PC,
			},
		}),
	)
}

func TestSession_RemoveCharacter(t *testing.T) {
	userID := id.New()
	s := New(id.NewURLFriendly(), userID, "test")
	character := s.AddCharacter(userID, PC, "test")

	ok := s.RemoveCharacter(character.ID)

	expect.That(t,
		is.EqualTo(ok, true),
		is.SliceOfLen(s.Characters, 0),
	)
}

func TestSession_FindCharacter(t *testing.T) {
	userID := id.New()
	s := New(id.NewURLFriendly(), userID, "test")
	c1 := s.AddCharacter(userID, PC, "test")
	c2 := s.AddCharacter(userID, PC, "test2")

	expect.That(t,
		is.DeepEqualTo(*c1, *s.FindCharacter(c1.ID)),
		is.DeepEqualTo(*c2, *s.FindCharacter(c2.ID)),
	)
}
