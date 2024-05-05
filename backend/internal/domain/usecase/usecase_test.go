package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
)

type repoMock struct {
	s session.Session
}

func (r *repoMock) Perform(ctx context.Context, id string, uow UnitOfWork) error {
	s, err := uow(ctx, r.s.ID == id, r.s)
	if errors.Is(err, NoSave) {
		return nil
	}

	if err != nil {
		return err
	}

	r.s = s

	return nil
}

type repoFixture struct {
	repo SessionRepository
}

func (f *repoFixture) BeforeEach(t *testing.TB) error {
	f.repo = &repoMock{}
	return nil
}

func TestLoadSession(t *testing.T) {
	repo := &repoMock{
		s: session.Session{
			ID:      "1",
			OwnerID: "2",
		},
	}
	loadSession := ProvideLoadSession(repo)

	t.Run("found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		got, err := loadSession(ctx, "1")
		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, session.Session{ID: "1", OwnerID: "2"}),
		)
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		got, err := loadSession(ctx, "2")

		expect.That(t,
			is.Error(err, ErrNotFound),
			is.DeepEqualTo(got, session.Session{}),
		)
	})

	t.Run("no_user", func(t *testing.T) {
		ctx := context.Background()
		got, err := loadSession(ctx, "1")

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.DeepEqualTo(got, session.Session{}),
		)
	})

	t.Run("not_authorized", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")
		got, err := loadSession(ctx, "1")

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.DeepEqualTo(got, session.Session{}),
		)
	})
}

func TestCreateSession(t *testing.T) {
	repo := &repoMock{}
	createSession := ProvideCreateSession(repo)

	t.Run("no_user", func(t *testing.T) {
		ctx := context.Background()
		got, err := createSession(ctx, CreateSessionRequest{})

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.DeepEqualTo(got, session.Session{}),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		got, err := createSession(ctx, CreateSessionRequest{
			Title: "Test",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(got, session.Session{OwnerID: "2", Title: "Test"}, is.ExcludeFields{"ID"}),
		)
	})
}

func TestJoinSession(t *testing.T) {
	repo := &repoMock{
		s: session.Session{
			ID:      "1",
			OwnerID: "2",
		},
	}
	joinSession := ProvideJoinSession(repo)

	t.Run("not_authorized", func(t *testing.T) {
		characterID, err := joinSession(context.Background(), JoinSessionRequest{
			SessionID:     "1",
			CharacterName: "Test",
		})

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(characterID, ""),
			is.SliceOfLen(repo.s.Characters, 0),
		)
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		characterID, err := joinSession(ctx, JoinSessionRequest{
			SessionID:     "99",
			CharacterName: "Test",
		})

		expect.That(t,
			is.Error(err, ErrNotFound),
			is.EqualTo(characterID, ""),
			is.SliceOfLen(repo.s.Characters, 0),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		characterID, err := joinSession(ctx, JoinSessionRequest{
			SessionID:     "1",
			CharacterName: "Test",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(repo.s, session.Session{
				ID:      "1",
				OwnerID: "2",
				Characters: []session.Character{
					{
						ID:      characterID,
						OwnerID: "3",
						Name:    "Test",
						Type:    session.PC,
					},
				},
			}),
		)
	})
}

func TestCreateAspect(t *testing.T) {
	repo := &repoMock{
		s: session.Session{
			ID:      "1",
			OwnerID: "2",
		},
	}
	createAspect := ProvideCreateAspect(repo)

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		aspectID, err := createAspect(ctx, CreateAspectRequest{
			SessionID: "2",
			Name:      "Test",
		})

		expect.That(t,
			is.Error(err, ErrNotFound),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("not_authorized", func(t *testing.T) {
		aspectID, err := createAspect(context.Background(), CreateAspectRequest{
			SessionID: "2",
			Name:      "Test",
		})

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("not_gm", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		aspectID, err := createAspect(ctx, CreateAspectRequest{
			SessionID: "1",
			Name:      "Test",
		})

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		aspectID, err := createAspect(ctx, CreateAspectRequest{
			SessionID: "1",
			Name:      "Test",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(repo.s.Aspects, []session.Aspect{
				{
					ID:   aspectID,
					Name: "Test",
				},
			}),
		)
	})
}

func TestCreateCharacterAspect(t *testing.T) {
	repo := &repoMock{
		s: session.Session{
			ID:      "1",
			OwnerID: "2",
			Characters: []session.Character{
				{
					ID:      "3",
					OwnerID: "4",
				},
			},
		},
	}
	createAspect := ProvideCreateCharacterAspect(repo)

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		aspectID, err := createAspect(ctx, CreateCharacterAspectRequest{
			CreateAspectRequest: CreateAspectRequest{
				SessionID: "2",
				Name:      "Test",
			},
			CharacterID: "3",
		})

		expect.That(t,
			is.Error(err, ErrNotFound),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("not_authorized", func(t *testing.T) {
		aspectID, err := createAspect(context.Background(), CreateCharacterAspectRequest{
			CreateAspectRequest: CreateAspectRequest{
				SessionID: "2",
				Name:      "Test",
			},
			CharacterID: "3",
		})

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("not_gm", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "3")

		aspectID, err := createAspect(ctx, CreateCharacterAspectRequest{
			CreateAspectRequest: CreateAspectRequest{
				SessionID: "1",
				Name:      "Test",
			},
			CharacterID: "3",
		})

		expect.That(t,
			is.Error(err, ErrForbidden),
			is.EqualTo(aspectID, ""),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("success", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")

		aspectID, err := createAspect(ctx, CreateCharacterAspectRequest{
			CreateAspectRequest: CreateAspectRequest{
				SessionID: "1",
				Name:      "Test",
			},
			CharacterID: "3",
		})

		expect.That(t,
			is.NoError(err),
			is.DeepEqualTo(repo.s.Characters[0].Aspects, []session.Aspect{
				{
					ID:   aspectID,
					Name: "Test",
				},
			}),
		)
	})
}

func TestDeleteAspect(t *testing.T) {
	repo := &repoMock{
		s: session.Session{
			ID:      "1",
			OwnerID: "2",
			Characters: []session.Character{
				{
					ID:      "3",
					OwnerID: "4",
					Aspects: session.Aspects{
						{ID: "5"},
					},
				},
			},
			Aspects: session.Aspects{
				{ID: "6"},
			},
		},
	}

	deleteAspect := ProvideDeleteAspect(repo)

	t.Run("not_authorized", func(t *testing.T) {
		err := deleteAspect(context.Background(), DeleteAspectRequest{
			SessionID: "1",
			AspectID:  "6",
		})
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, DeleteAspectRequest{
			SessionID: "2",
			AspectID:  "6",
		})
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("not_gm", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "4")
		err := deleteAspect(ctx, DeleteAspectRequest{
			SessionID: "1",
			AspectID:  "6",
		})
		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("aspect_not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, DeleteAspectRequest{
			SessionID: "1",
			AspectID:  "7",
		})
		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("success_global_aspect", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, DeleteAspectRequest{
			SessionID: "1",
			AspectID:  "6",
		})
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(repo.s.Aspects, 0),
		)
	})

	t.Run("success_character_aspect", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := deleteAspect(ctx, DeleteAspectRequest{
			SessionID: "1",
			AspectID:  "5",
		})
		expect.That(t,
			is.NoError(err),
			is.SliceOfLen(repo.s.Characters[0].Aspects, 0),
		)
	})
}

func TestUpdateFatePoints(t *testing.T) {
	repo := &repoMock{
		s: session.Session{
			ID:      "1",
			OwnerID: "2",
			Characters: []session.Character{
				{
					ID:      "3",
					OwnerID: "4",
				},
			},
		},
	}
	updateFatePoints := ProvideUpdateFatePoints(repo)

	t.Run("not_authorized", func(t *testing.T) {
		err := updateFatePoints(context.Background(), UpdateFatePointsRequest{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := updateFatePoints(ctx, UpdateFatePointsRequest{
			SessionID:   "99",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, ErrNotFound))
	})

	t.Run("character_not_found", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := updateFatePoints(ctx, UpdateFatePointsRequest{
			SessionID:   "1",
			CharacterID: "4",
			Delta:       1,
		})

		expect.That(t, is.Error(err, ErrInvalidCharacter))
	})

	t.Run("neither_gm_nor_player", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "5")
		err := updateFatePoints(ctx, UpdateFatePointsRequest{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("player_but_wants_to_increase", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "4")
		err := updateFatePoints(ctx, UpdateFatePointsRequest{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t, is.Error(err, ErrForbidden))
	})

	t.Run("gm_increase", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "2")
		err := updateFatePoints(ctx, UpdateFatePointsRequest{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       1,
		})

		expect.That(t,
			is.NoError(err),
			is.EqualTo(repo.s.Characters[0].FatePoints, 1),
		)
	})

	t.Run("player_decrease", func(t *testing.T) {
		ctx := auth.WithUserID(context.Background(), "4")
		err := updateFatePoints(ctx, UpdateFatePointsRequest{
			SessionID:   "1",
			CharacterID: "3",
			Delta:       -1,
		})

		expect.That(t,
			is.NoError(err),
			is.EqualTo(repo.s.Characters[0].FatePoints, 0),
		)
	})

}
