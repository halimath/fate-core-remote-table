//go:generate stringer -type=CharacterType -output charactertype_gen.go
package session

import (
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
)

type CharacterType int

const (
	PC CharacterType = iota
	NPC
)

type entity interface {
	id() id.ID
}

func removeByID[T entity](e *[]T, entityID id.ID) bool {
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
	ID   id.ID
	Name string
}

func (a Aspect) id() id.ID {
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

func (a *Aspects) RemoveAspect(aspectID id.ID) bool {
	return removeByID((*[]Aspect)(a), aspectID)
}

type Character struct {
	ID         id.ID
	OwnerID    id.ID
	Type       CharacterType
	Name       string
	FatePoints int
	Aspects
}

func (c Character) id() id.ID {
	return c.ID
}

type Session struct {
	ID           id.ID
	LastModified time.Time
	OwnerID      id.ID
	Title        string
	Characters   []Character
	Aspects
}

func New(sessionID, ownerID id.ID, title string) Session {
	return Session{
		ID:           sessionID,
		LastModified: time.Now().Truncate(time.Millisecond),
		OwnerID:      ownerID,
		Title:        title,
		Characters:   make([]Character, 0),
		Aspects:      make([]Aspect, 0),
	}
}

func (s *Session) AddCharacter(ownerID id.ID, typ CharacterType, name string, aspects ...Aspect) *Character {
	s.Characters = append(s.Characters, Character{
		ID:      id.New(),
		OwnerID: ownerID,
		Type:    typ,
		Name:    name,
		Aspects: aspects,
	})

	return &(s.Characters[len(s.Characters)-1])
}

func (s *Session) RemoveCharacter(characterID id.ID) bool {
	return removeByID(&s.Characters, characterID)
}

func (s *Session) FindCharacter(characterID id.ID) *Character {
	for i := range s.Characters {
		if s.Characters[i].ID == characterID {
			return &s.Characters[i]
		}
	}

	return nil
}
