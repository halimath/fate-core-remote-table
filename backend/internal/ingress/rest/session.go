//go:generate oapi-codegen -package rest -generate types -o dtos_gen.go ../../../../docs/api.yaml
package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/domain"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createaspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createcharacteraspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/deleteaspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/joinsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/loadsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/updatefatepoints"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/httputils/errmux"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func newSessionAPIHandler(
	cfg config.Config,
	createSession createsession.Port,
	loadSession loadsession.Port,
	joinSession joinsession.Port,
	createAspect createaspect.Port,
	createCharacterAspect createcharacteraspect.Port,
	deleteAspect deleteaspect.Port,
	updateFatePoints updatefatepoints.Port,
) http.Handler {
	mux := errmux.NewServeMux()
	mux.ErrorHandler = handleError

	mux.Handle("POST /", createSessionHandler(createSession))
	mux.Handle("GET /{id}", getSessionHandler(cfg, loadSession))
	mux.Handle("POST /{id}/join", joinSessionHandler(joinSession))
	mux.Handle("POST /{id}/aspects", createAspectHandler(createAspect))
	mux.Handle("POST /{id}/characters/{characterID}/aspects", createCharacterAspectHandler(createCharacterAspect))
	mux.Handle("DELETE /{id}/aspects/{aspectID}", deleteAspectHandler(deleteAspect))
	// mux.HandleFunc("POST /{id}/characters", wrapper.CreateCharacter)
	// mux.HandleFunc("DELETE /{id}/characters/{characterId}", wrapper.DeleteCharacter)
	mux.Handle("PUT /{id}/characters/{characterID}/fatepoints", updateFatePointsHandler(updateFatePoints))

	return mux
}

func updateFatePointsHandler(updateFatePoints updatefatepoints.Port) errmux.Handler {
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

		err := updateFatePoints(r.Context(), updatefatepoints.Request{
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

func deleteAspectHandler(deleteAspect deleteaspect.Port) errmux.Handler {
	return errmux.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		err := deleteAspect(r.Context(), deleteaspect.Request{
			SessionID: r.PathValue("id"),
			AspectID:  r.PathValue("aspectID"),
		})

		if err != nil {
			return err
		}

		return response.NoContent(w, r, response.StatusCode(http.StatusAccepted))
	})
}

func createCharacterAspectHandler(createCharacterAspect createcharacteraspect.Port) errmux.Handler {
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

		aspectID, err := createCharacterAspect(r.Context(), createcharacteraspect.Request{
			SessionID:   r.PathValue("id"),
			Name:        body.Name,
			CharacterID: r.PathValue("characterID"),
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, aspectID)

	})
}

func createAspectHandler(createAspect createaspect.Port) errmux.Handler {
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

		aspectID, err := createAspect(r.Context(), createaspect.Request{
			SessionID: r.PathValue("id"),
			Name:      body.Name,
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, aspectID)

	})
}

func joinSessionHandler(joinSession joinsession.Port) errmux.Handler {
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

		characterID, err := joinSession(r.Context(), joinsession.Request{
			SessionID:     r.PathValue("id"),
			CharacterName: body.Name,
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, characterID)
	})
}

func createSessionHandler(createSession createsession.Port) errmux.Handler {
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

		ses, err := createSession(r.Context(), createsession.Request{
			Title: body.Title,
		})

		if err != nil {
			return err
		}

		return response.PlainText(w, r, ses.ID, response.StatusCode(http.StatusCreated))
	})
}

func getSessionHandler(cfg config.Config, loadSession loadsession.Port) errmux.Handler {
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

		// By default, allow session responses to be cache, but only in browsers (private) and force browser
		// to refresh cache on any request (no-cache).
		cacheHeaderOption := response.AddHeader("Cache-Control", "private, no-cache")
		if cfg.DevMode {
			// In dev mode, disable the cache (no-store)
			cacheHeaderOption = response.AddHeader("Cache-Control", "no-store")
		}

		return response.JSON(w, r, convertSession(ses),
			response.AddHeader("Last-Modified", ses.LastModified.UTC().Truncate(time.Second).Format(http.TimeFormat)),
			cacheHeaderOption,
		)
	})
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, domain.ErrNotFound) {
		response.NotFound(w, r)
		return
	}

	if errors.Is(err, domain.ErrForbidden) {
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

func bindBody(r *http.Request, payload any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, payload)
}
