//go:generate oapi-codegen -package web -generate types -o dtos_gen.go ../../../docs/api.yaml
package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecase"
	"github.com/halimath/httputils/errmux"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func newSessionAPIHandler(
	createSession usecase.CreateSession,
	loadSession usecase.LoadSession,
	joinSession usecase.JoinSession,
	createAspect usecase.CreateAspect,
	createCharacterAspect usecase.CreateCharacterAspect,
	deleteAspect usecase.DeleteAspect,
	updateFatePoints usecase.UpdateFatePoints,
) http.Handler {
	mux := errmux.NewServeMux()
	mux.ErrorHandler = handleError

	mux.Handle("POST /", createSessionHandler(createSession))
	mux.Handle("GET /{id}", getSessionHandler(loadSession))
	mux.Handle("POST /{id}/join", joinSessionHandler(joinSession))
	mux.Handle("POST /{id}/aspects", createAspectHandler(createAspect))
	mux.Handle("POST /{id}/characters/{characterID}/aspects", createCharacterAspectHandler(createCharacterAspect))
	mux.Handle("DELETE /{id}/aspects/{aspectID}", deleteAspectHandler(deleteAspect))
	// mux.HandleFunc("POST /{id}/characters", wrapper.CreateCharacter)
	// mux.HandleFunc("DELETE /{id}/characters/{characterId}", wrapper.DeleteCharacter)
	mux.Handle("PUT /{id}/characters/{characterID}/fatepoints", updateFatePointsHandler(updateFatePoints))

	return mux
}

func updateFatePointsHandler(updateFatePoints usecase.UpdateFatePoints) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var body UpdateFatePoints

		if err := bindBody(r, &body); err != nil {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidUpdateFatePoints",
				Title:  "Invalid request payload to update fate points",
				Status: http.StatusBadRequest,
				Errors: []any{err},
			})
		}

		err := updateFatePoints(r.Context(), usecase.UpdateFatePointsRequest{
			SessionID:   r.PathValue("id"),
			CharacterID: r.PathValue("characterID"),
			Delta:       body.FatePointsDelta,
		})

		if err != nil {
			return err
		}

		return response.NoContent(w, r)
	})
}

func deleteAspectHandler(deleteAspect usecase.DeleteAspect) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		err := deleteAspect(r.Context(), usecase.DeleteAspectRequest{
			SessionID: r.PathValue("id"),
			AspectID:  r.PathValue("aspectID"),
		})

		if err != nil {
			return err
		}

		return response.NoContent(w, r, response.StatusCode(http.StatusAccepted))
	})
}

func createCharacterAspectHandler(createCharacterAspect usecase.CreateCharacterAspect) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var body CreateAspect

		if err := bindBody(r, &body); err != nil {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidCreateAspect",
				Title:  "Invalid request payload to create aspect",
				Status: http.StatusBadRequest,
				Errors: []any{err},
			})
		}

		if body.Name == "" {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidCreateAspect",
				Title:  "Missing aspect name",
				Status: http.StatusBadRequest,
			})
		}

		aspectID, err := createCharacterAspect(r.Context(), usecase.CreateCharacterAspectRequest{
			CreateAspectRequest: usecase.CreateAspectRequest{
				SessionID: r.PathValue("id"),
				Name:      body.Name,
			},
			CharacterID: r.PathValue("characterID"),
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, aspectID)

	})
}

func createAspectHandler(createAspect usecase.CreateAspect) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var body CreateAspect

		if err := bindBody(r, &body); err != nil {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidCreateAspect",
				Title:  "Invalid request payload to create aspect",
				Status: http.StatusBadRequest,
				Errors: []any{err},
			})
		}

		if body.Name == "" {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidCreateAspect",
				Title:  "Missing aspect name",
				Status: http.StatusBadRequest,
			})
		}

		aspectID, err := createAspect(r.Context(), usecase.CreateAspectRequest{
			SessionID: r.PathValue("id"),
			Name:      body.Name,
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, aspectID)

	})
}

func joinSessionHandler(joinSession usecase.JoinSession) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var body JoinSession
		if err := bindBody(r, &body); err != nil {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidJoinSession",
				Title:  "Invalid request payload to join a session",
				Status: http.StatusBadRequest,
				Errors: []any{err},
			})
		}

		if body.Name == "" {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidJoinSession",
				Title:  "Missing character name",
				Status: http.StatusBadRequest,
			})
		}

		characterID, err := joinSession(r.Context(), usecase.JoinSessionRequest{
			SessionID:     r.PathValue("id"),
			CharacterName: body.Name,
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, characterID)
	})
}

func createSessionHandler(createSession usecase.CreateSession) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		var body CreateSession
		if err := bindBody(r, &body); err != nil {
			return response.Problem(w, r, response.ProblemDetails{
				Type:   "github.com/halimath/fate-table/problem/invalidSessionCreate",
				Title:  "Invalid session creation payload",
				Status: http.StatusBadRequest,
				Errors: []any{err},
			})
		}

		ses, err := createSession(r.Context(), usecase.CreateSessionRequest{
			Title: body.Title,
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, ses.ID, response.StatusCode(http.StatusCreated))
	})
}

func getSessionHandler(loadSession usecase.LoadSession) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		sessionID := r.PathValue("id")
		ses, err := loadSession(r.Context(), sessionID)
		if err != nil {
			return err
		}

		ifModifiedSince := r.Header.Get("If-Modified-Since")
		if len(ifModifiedSince) > 0 {
			ifModifiedSinceTime, err := http.ParseTime(ifModifiedSince)
			if err != nil {
				kvlog.FromContext(r.Context()).Logs("failed to parse If-Modified-Since header", kvlog.WithErr(err))
			} else {
				if !ses.LastModified.UTC().Truncate(time.Second).After(ifModifiedSinceTime.UTC().Truncate(time.Second)) {
					return response.NotModified(w, r)
				}
			}
		}

		return response.JSON(w, r, convertSession(ses),
			response.AddHeader("Last-Modified", ses.LastModified.UTC().Truncate(time.Second).Format(http.TimeFormat)),
			response.AddHeader("Cache-Control", "private, no-cache"),
		)
	})
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, usecase.ErrNotFound) {
		response.NotFound(w, r)
		return
	}

	if errors.Is(err, usecase.ErrForbidden) {
		response.Forbidden(w, r)
		return
	}

	response.Error(w, r, err)
}

// func convertCharacterType(t CreateCharacterType) (session.CharacterType, error) {
// 	switch t {
// 	case "PC":
// 		return session.PC, nil
// 	case "NPC":
// 		return session.NPC, nil
// 	default:
// 		return 0, fmt.Errorf("invalid character type: %s", t)
// 	}
// }

func convertSession(s session.Session) Session {
	return Session{
		Id:         s.ID,
		OwnerId:    s.OwnerID,
		Title:      s.Title,
		Aspects:    convertAspects(s.Aspects),
		Characters: convertCharacters(s.Characters),
	}
}

func convertCharacters(cs []session.Character) []Character {
	if len(cs) == 0 {
		return []Character{}
	}

	res := make([]Character, len(cs))

	for i, c := range cs {
		res[i] = Character{
			Name:       c.Name,
			Type:       CharacterType(c.Type.String()),
			Id:         c.ID,
			OwnerId:    c.OwnerID,
			FatePoints: c.FatePoints,
			Aspects:    convertAspects(c.Aspects),
		}
	}

	return res
}

func convertAspects(a []session.Aspect) []Aspect {
	if len(a) == 0 {
		return []Aspect{}
	}

	res := make([]Aspect, len(a))

	for i, aspect := range a {
		res[i] = Aspect{
			Id:   aspect.ID,
			Name: aspect.Name,
		}
	}

	return res
}
