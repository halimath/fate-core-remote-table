//go:generate oapi-codegen -package boundary -generate types,std-http -o rest_gen.go ../../../docs/api.yaml

package boundary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/boundary/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/control"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/session"
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
	controller   control.SessionController
	authProvider auth.Provider
}

var _ ServerInterface = &restHandler{}

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

	existingToken, ok := auth.ExtractBearerToken(r)
	if ok {
		kvlog.L.Logs("renewToken")
		token, err = h.authProvider.RenewToken(existingToken)
	} else {
		kvlog.L.Logs("createToken")
		token, err = h.authProvider.CreateToken()
	}

	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.PlainText(w, r, token, response.StatusCode(http.StatusCreated))
}

func (h *restHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateSession

	if err := bindBody(r, &dto); err != nil {
		response.Error(w, r, err)
		return
	}

	sessionID, err := h.controller.Create(r.Context(), userID, dto.Title)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.PlainText(w, r, sessionID.String(), response.StatusCode(http.StatusCreated))
}

func (h *restHandler) GetSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	s, err := h.controller.Load(r.Context(), id.FromString(sessionID))
	if err != nil {
		response.Error(w, r, err)
		return
	}

	ifModifiedSince := r.Header.Get("If-Modified-Since")
	if ifModifiedSince != "" {
		cacheDate, err := time.Parse(time.RFC1123, ifModifiedSince)
		if err == nil {
			if !cacheDate.UTC().Truncate(time.Second).Before(s.LastModified.UTC().Truncate(time.Second)) {
				response.NotModified(w, r)
				return
			}
		}
	}

	header := w.Header()
	header.Add("Last-Modified", s.LastModified.In(GMT).Format(time.RFC1123))
	header.Add("Cache-Control", "private, no-cache")

	response.JSON(w, r, convertSession(s))
}

func (h *restHandler) CreateAspect(w http.ResponseWriter, r *http.Request, sessionID string) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateAspect

	if err := bindBody(r, &dto); err != nil {
		response.Error(w, r, err)
		return
	}

	aspectID, err := h.controller.CreateAspect(r.Context(), userID, id.FromString(sessionID), dto.Name)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.PlainText(w, r, aspectID.String(), response.StatusCode(http.StatusCreated))
}

func (h *restHandler) DeleteAspect(w http.ResponseWriter, r *http.Request, sessionID string, aspectID string) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	err := h.controller.DeleteAspect(r.Context(), userID, id.FromString(sessionID), id.FromString(aspectID))
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.NoContent(w, r)
}

func (h *restHandler) CreateCharacter(w http.ResponseWriter, r *http.Request, sessionID string) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateCharacter

	if err := bindBody(r, &dto); err != nil {
		response.Error(w, r, err)
		return
	}

	t, err := convertCharacterType(dto.Type)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	characterID, err := h.controller.CreateCharacter(r.Context(), userID, id.FromString(sessionID), t, dto.Name)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.PlainText(w, r, characterID.String(), response.StatusCode(http.StatusCreated))
}

func (h *restHandler) DeleteCharacter(w http.ResponseWriter, r *http.Request, sessionID, characterID string) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	err := h.controller.DeleteCharacter(r.Context(), userID, id.FromString(sessionID), id.FromString(characterID))
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.NoContent(w, r)
}

func (h *restHandler) CreateCharacterAspect(w http.ResponseWriter, r *http.Request, sessionID, characterID string) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto CreateAspect

	if err := bindBody(r, &dto); err != nil {
		response.Error(w, r, err)
		return
	}

	aspectID, err := h.controller.CreateCharacterAspect(r.Context(), userID, id.FromString(sessionID), id.FromString(characterID), dto.Name)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	response.PlainText(w, r, aspectID.String(), response.StatusCode(http.StatusCreated))
}

func (h *restHandler) UpdateFatePoints(w http.ResponseWriter, r *http.Request, sessionID, characterID string) {
	userID, ok := auth.UserID(r)
	if !ok {
		response.Forbidden(w, r)
		return
	}

	var dto UpdateFatePoints

	if err := bindBody(r, &dto); err != nil {
		response.Error(w, r, err)
		return

	}

	err := h.controller.UpdateFatePoints(r.Context(), userID, id.FromString(sessionID), id.FromString(characterID), dto.FatePointsDelta)
	if err != nil {
		response.Error(w, r, err)
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
		Id:         s.ID.String(),
		OwnerId:    s.OwnerID.String(),
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
			Id:         c.ID.String(),
			OwnerId:    c.OwnerID.String(),
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
			Id:   aspect.ID.String(),
			Name: aspect.Name,
		}
	}

	return res
}
