//go:generate oapi-codegen -package boundary -generate types,std-http -o rest_gen.go ../../../docs/api.yaml

package boundary

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

var (
	GMT *time.Location
)

func init() {
	var err error
	GMT, err = time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
}

type restHandler struct {
	versionInfo  VersionInfo
	service      session.Service
	authProvider *auth.Manager
}

var _ ServerInterface = &restHandler{}

type HTTPError interface {
	error
	StatusCode() int
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	e := Error{
		Error: err.Error(),
		Code:  http.StatusInternalServerError,
	}

	if httpError, ok := err.(HTTPError); ok {
		e.Code = httpError.StatusCode()
	} else if errors.Is(err, session.ErrNotFound) {
		e.Code = http.StatusNotFound
	}

	response.JSON(w, r, e, response.StatusCode(e.Code))
}

func bindBody(r *http.Request, payload any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, payload)
}

func (h *restHandler) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var token string
	var err error

	existingToken, ok := extractBearerToken(r)
	if ok {
		kvlog.L.Logs("renewToken")
		token, err = h.authProvider.RenewToken(existingToken)
	} else {
		kvlog.L.Logs("createToken")
		token, err = h.authProvider.CreateToken()
	}

	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error creating auth token", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.PlainText(w, r, token, response.StatusCode(http.StatusCreated))
}

func (h *restHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateSession

	if err := bindBody(r, &dto); err != nil {
		kvlog.FromContext(r.Context()).Logs("error unmarshaling request payload", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	sessionID, err := h.service.Create(r.Context(), userID, dto.Title)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error creating session", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.PlainText(w, r, sessionID, response.StatusCode(http.StatusCreated))
}

func (h *restHandler) GetSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	s, err := h.service.Load(r.Context(), sessionID)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error loading session", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if ifModifiedSince != "" {
		cacheDate, err := http.ParseTime(ifModifiedSince)
		if err == nil {
			if cacheDate.UTC().Truncate(time.Second).After(s.LastModified.UTC().Truncate(time.Second)) {
				response.NotModified(w, r)
				return
			}
		}
	}

	response.JSON(w, r, convertSession(s),
		response.SetHeader("Last-Modified", s.LastModified.In(GMT).Format(time.RFC1123), true),
		response.SetHeader("Cache-Control", "private, must-revalidate", true),
	)
}

func (h *restHandler) CreateAspect(w http.ResponseWriter, r *http.Request, sessionID string) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateAspect

	if err := bindBody(r, &dto); err != nil {
		kvlog.FromContext(r.Context()).Logs("error unmarshaling request payload", kvlog.WithErr(err))
		response.Error(w, r, err, response.StatusCode(http.StatusBadRequest))
		return
	}

	aspectID, err := h.service.CreateAspect(r.Context(), userID, sessionID, dto.Name)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error creating aspect", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.PlainText(w, r, aspectID, response.StatusCode(http.StatusCreated))
}

func (h *restHandler) DeleteAspect(w http.ResponseWriter, r *http.Request, sessionID string, aspectID string) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	err := h.service.DeleteAspect(r.Context(), userID, sessionID, aspectID)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error deleting aspect", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.NoContent(w, r)
}

func (h *restHandler) CreateCharacter(w http.ResponseWriter, r *http.Request, sessionID string) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateCharacter

	if err := bindBody(r, &dto); err != nil {
		kvlog.FromContext(r.Context()).Logs("error unmarshaling request payload", kvlog.WithErr(err))
		response.Error(w, r, err, response.StatusCode(http.StatusBadRequest))
		return
	}

	t, err := convertCharacterType(dto.Type)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error converting character type", kvlog.WithErr(err))
		response.Error(w, r, err, response.StatusCode(http.StatusBadRequest))
		return
	}

	characterID, err := h.service.CreateCharacter(r.Context(), userID, sessionID, t, dto.Name)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error creating character", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.PlainText(w, r, characterID, response.StatusCode(http.StatusCreated))
}

func (h *restHandler) DeleteCharacter(w http.ResponseWriter, r *http.Request, sessionID, characterID string) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	err := h.service.DeleteCharacter(r.Context(), userID, sessionID, characterID)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error deleting character", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.NoContent(w, r)
}

func (h *restHandler) CreateCharacterAspect(w http.ResponseWriter, r *http.Request, sessionID, characterID string) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateAspect

	if err := bindBody(r, &dto); err != nil {
		kvlog.FromContext(r.Context()).Logs("error unmarshaling request payload", kvlog.WithErr(err))
		response.Error(w, r, err, response.StatusCode(http.StatusBadRequest))
		return
	}

	aspectID, err := h.service.CreateCharacterAspect(r.Context(), userID, sessionID, characterID, dto.Name)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error creating character aspect", kvlog.WithErr(err))
		handleError(w, r, err)
		return
	}

	response.PlainText(w, r, aspectID, response.StatusCode(http.StatusCreated))
}

func (h *restHandler) UpdateFatePoints(w http.ResponseWriter, r *http.Request, sessionID, characterID string) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto UpdateFatePoints

	if err := bindBody(r, &dto); err != nil {
		kvlog.FromContext(r.Context()).Logs("error unmarshaling request payload", kvlog.WithErr(err))
		response.Error(w, r, err, response.StatusCode(http.StatusBadRequest))
		return

	}

	err := h.service.UpdateFatePoints(r.Context(), userID, sessionID, characterID, dto.FatePointsDelta)
	if err != nil {
		kvlog.FromContext(r.Context()).Logs("error updating fate points", kvlog.WithErr(err))
		handleError(w, r, err)
		return

	}

	response.NoContent(w, r)
}

func (h *restHandler) GetVersionInfo(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, r, h.versionInfo)
}

func convertCharacterType(t CreateCharacterType) (session.CharacterType, error) {
	switch t {
	case "PC":
		return session.PC, nil
	case "NPC":
		return session.NPC, nil
	default:
		return 0, fmt.Errorf("invalid character type: %s", t)
	}
}

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
