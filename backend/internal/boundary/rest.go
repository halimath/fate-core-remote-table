//go:generate oapi-codegen -generate types,server -package boundary -o rest_gen.go ../../../docs/api.yaml

package boundary

import (
	"fmt"
	"net/http"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/boundary/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/control"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/session"
	"github.com/halimath/kvlog"
	"github.com/labstack/echo/v4"
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

func (h *restHandler) CreateAuthToken(ctx echo.Context) error {
	if auth.IsAuthorized(ctx) {
		kvlog.Warn(kvlog.Evt("alreadyAuthorized"))
		return echo.ErrForbidden
	}

	token, err := h.authProvider.CreateToken()
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusCreated, "text/plain", []byte(token))
}

func (h *restHandler) CreateSession(ctx echo.Context) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	var dto CreateSession

	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	sessionID, err := h.controller.Create(ctx.Request().Context(), userID, dto.Title)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusCreated, "text/plain", []byte(sessionID))
}

func (h *restHandler) GetSession(ctx echo.Context, sessionID string) error {
	s, err := h.controller.Load(ctx.Request().Context(), id.FromString(sessionID))
	if err != nil {
		return err
	}

	ifModifiedSince := ctx.Request().Header.Get("If-Modified-Since")
	if ifModifiedSince != "" {
		cacheDate, err := time.Parse(time.RFC1123, ifModifiedSince)
		if err == nil {
			if !cacheDate.UTC().Truncate(time.Second).Before(s.LastModified.UTC().Truncate(time.Second)) {
				return ctx.NoContent(http.StatusNotModified)
			}
		}
	}

	header := ctx.Response().Header()
	header.Add("Last-Modified", s.LastModified.In(GMT).Format(time.RFC1123))
	header.Add("Cache-Control", "private; must-revalidate")

	return ctx.JSON(http.StatusOK, convertSession(s))
}

func (h *restHandler) CreateAspect(ctx echo.Context, sessionID string) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	var dto CreateAspect

	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	aspectID, err := h.controller.CreateAspect(ctx.Request().Context(), userID, id.FromString(sessionID), dto.Name)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusCreated, "text/plain", []byte(aspectID))
}

func (h *restHandler) DeleteAspect(ctx echo.Context, sessionID string, aspectID string) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	err := h.controller.DeleteAspect(ctx.Request().Context(), userID, id.FromString(sessionID), id.FromString(aspectID))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *restHandler) CreateCharacter(ctx echo.Context, sessionID string) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	var dto CreateCharacter

	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	t, err := convertCharacterType(dto.Type)
	if err != nil {
		return err
	}

	characterID, err := h.controller.CreateCharacter(ctx.Request().Context(), userID, id.FromString(sessionID), t, dto.Name)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusCreated, "text/plain", []byte(characterID))
}

func (h *restHandler) DeleteCharacter(ctx echo.Context, sessionID, characterID string) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	err := h.controller.DeleteCharacter(ctx.Request().Context(), userID, id.FromString(sessionID), id.FromString(characterID))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *restHandler) CreateCharacterAspect(ctx echo.Context, sessionID, characterID string) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	var dto CreateAspect

	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	aspectID, err := h.controller.CreateCharacterAspect(ctx.Request().Context(), userID, id.FromString(sessionID), id.FromString(characterID), dto.Name)
	if err != nil {
		return err
	}

	return ctx.Blob(http.StatusCreated, "text/plain", []byte(aspectID))
}

func (h *restHandler) UpdateFatePoints(ctx echo.Context, sessionID, characterID string) error {
	userID, ok := auth.UserID(ctx)
	if !ok {
		return echo.ErrForbidden
	}

	var dto UpdateFatePoints

	if err := ctx.Bind(&dto); err != nil {
		return err
	}

	err := h.controller.UpdateFatePoints(ctx.Request().Context(), userID, id.FromString(sessionID), id.FromString(characterID), dto.FatePointsDelta)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *restHandler) GetVersionInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, h.versionInfo)
}

func convertCharacterType(t string) (session.CharacterType, error) {
	switch t {
	case "PC":
		return session.PC, nil
	case "NPC":
		return session.NPC, nil
	default:
		return 0, fmt.Errorf("%w: invalid character type: %s", echo.ErrBadRequest, t)
	}
}

func convertSession(s session.Session) Session {
	return Session{
		Id:      s.ID.String(),
		OwnerId: s.OwnerID.String(),
		CreateSession: CreateSession{
			Title: s.Title,
		},
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
			CreateCharacter: CreateCharacter{
				Name: c.Name,
				Type: c.Type.String(),
			},
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
			Id: aspect.ID.String(),
			CreateAspect: CreateAspect{
				Name: aspect.Name,
			},
		}
	}

	return res
}
