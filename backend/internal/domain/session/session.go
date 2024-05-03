//go:generate stringer -type=CharacterType -output charactertype_gen.go
package session

import (
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/id"
)

type CharacterType int

const (
	PC CharacterType = iota
	NPC
)

type entity interface {
	id() string
}

func removeByID[T entity](e *[]T, entityID string) bool {
	for i := range *e {
		if (*e)[i].id() == entityID {
			temp := (*e)[:i]
			temp = append(temp, (*e)[i+1:]...)
			*e = temp

			return true
		}
	}

	return false
}

type Aspect struct {
	ID   string
	Name string
}

func (a Aspect) id() string {
	return a.ID
}

type Aspects []Aspect

func (a *Aspects) AddAspect(name string) *Aspect {
	*a = append(*a, Aspect{
		ID:   id.New(),
		Name: name,
	})
	return &([]Aspect(*a)[len(*a)-1])
}

func (a *Aspects) RemoveAspect(aspectID string) bool {
	return removeByID((*[]Aspect)(a), aspectID)
}

type Character struct {
	ID         string
	OwnerID    string
	Type       CharacterType
	Name       string
	FatePoints int
	Aspects
}

func (c Character) id() string {
	return c.ID
}

type Session struct {
	ID           string
	LastModified time.Time
	OwnerID      string
	Title        string
	Characters   []Character
	Aspects
}

func New(sessionID, ownerID string, title string) Session {
	return Session{
		ID:           sessionID,
		LastModified: time.Now().Truncate(time.Millisecond),
		OwnerID:      ownerID,
		Title:        title,
		Characters:   make([]Character, 0),
		Aspects:      make([]Aspect, 0),
	}
}

func (s *Session) IsMember(userID string) bool {
	if s.OwnerID == userID {
		return true
	}

	for _, c := range s.Characters {
		if c.OwnerID == userID {
			return true
		}
	}

	return false
}

func (s *Session) AddCharacter(ownerID string, typ CharacterType, name string, aspects ...Aspect) *Character {
	s.Characters = append(s.Characters, Character{
		ID:      id.New(),
		OwnerID: ownerID,
		Type:    typ,
		Name:    name,
		Aspects: aspects,
	})

	return &(s.Characters[len(s.Characters)-1])
}

func (s *Session) RemoveCharacter(characterID string) bool {
	return removeByID(&s.Characters, characterID)
}

func (s *Session) FindCharacter(characterID string) *Character {
	for i := range s.Characters {
		if s.Characters[i].ID == characterID {
			return &s.Characters[i]
		}
	}

	return nil
}
