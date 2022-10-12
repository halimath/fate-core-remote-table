// Package boundary provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package boundary

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

const (
	BearerScopes = "bearer.Scopes"
)

// Defines values for CharacterType.
const (
	CharacterTypeNPC CharacterType = "NPC"
	CharacterTypePC  CharacterType = "PC"
)

// Defines values for CreateCharacterType.
const (
	CreateCharacterTypeNPC CreateCharacterType = "NPC"
	CreateCharacterTypePC  CreateCharacterType = "PC"
)

// Aspect defines model for Aspect.
type Aspect struct {
	// The unique id of the aspect
	Id string `json:"id"`

	// The aspect's name
	Name string `json:"name"`
}

// Character defines model for Character.
type Character struct {
	Aspects []Aspect `json:"aspects"`

	// Non-negative number of Fate Points for the character
	FatePoints int `json:"fatePoints"`

	// The unique id of the character
	Id string `json:"id"`

	// The character's name
	Name string `json:"name"`

	// The unique id of the characters's owner
	OwnerId string        `json:"ownerId"`
	Type    CharacterType `json:"type"`
}

// CharacterType defines model for Character.Type.
type CharacterType string

// CreateAspect defines model for CreateAspect.
type CreateAspect struct {
	// The aspect's name
	Name string `json:"name"`
}

// CreateCharacter defines model for CreateCharacter.
type CreateCharacter struct {
	// The character's name
	Name string              `json:"name"`
	Type CreateCharacterType `json:"type"`
}

// CreateCharacterType defines model for CreateCharacter.Type.
type CreateCharacterType string

// CreateSession defines model for CreateSession.
type CreateSession struct {
	// Human readable title of the session
	Title string `json:"title"`
}

// Error defines model for Error.
type Error struct {
	// error code
	Code int `json:"code"`

	// Human-readable error message
	Error string `json:"error"`
}

// Session defines model for Session.
type Session struct {
	Aspects    []Aspect    `json:"aspects"`
	Characters []Character `json:"characters"`

	// The unique id of the session
	Id string `json:"id"`

	// The unique id of the session's owner
	OwnerId string `json:"ownerId"`

	// Human readable title of the session
	Title string `json:"title"`
}

// UpdateFatePoints defines model for UpdateFatePoints.
type UpdateFatePoints struct {
	// Number to modify character's Fate Points (negative or positive)
	FatePointsDelta int `json:"fatePointsDelta"`
}

// VersionInfo defines model for VersionInfo.
type VersionInfo struct {
	// The version string of the API specs.
	ApiVersion string `json:"apiVersion"`

	// Git commit hash of the backend code.
	Commit string `json:"commit"`

	// The version string of the backend component.
	Version string `json:"version"`
}

// CreateSessionJSONBody defines parameters for CreateSession.
type CreateSessionJSONBody = CreateSession

// CreateAspectJSONBody defines parameters for CreateAspect.
type CreateAspectJSONBody = CreateAspect

// CreateCharacterJSONBody defines parameters for CreateCharacter.
type CreateCharacterJSONBody = CreateCharacter

// CreateCharacterAspectJSONBody defines parameters for CreateCharacterAspect.
type CreateCharacterAspectJSONBody = CreateAspect

// UpdateFatePointsJSONBody defines parameters for UpdateFatePoints.
type UpdateFatePointsJSONBody = UpdateFatePoints

// CreateSessionJSONRequestBody defines body for CreateSession for application/json ContentType.
type CreateSessionJSONRequestBody = CreateSessionJSONBody

// CreateAspectJSONRequestBody defines body for CreateAspect for application/json ContentType.
type CreateAspectJSONRequestBody = CreateAspectJSONBody

// CreateCharacterJSONRequestBody defines body for CreateCharacter for application/json ContentType.
type CreateCharacterJSONRequestBody = CreateCharacterJSONBody

// CreateCharacterAspectJSONRequestBody defines body for CreateCharacterAspect for application/json ContentType.
type CreateCharacterAspectJSONRequestBody = CreateCharacterAspectJSONBody

// UpdateFatePointsJSONRequestBody defines body for UpdateFatePoints for application/json ContentType.
type UpdateFatePointsJSONRequestBody = UpdateFatePointsJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create an authorization token for the client
	// (POST /auth/new)
	CreateAuthToken(ctx echo.Context) error
	// Create a new session
	// (POST /sessions)
	CreateSession(ctx echo.Context) error
	// Get the session with the given id
	// (GET /sessions/{id})
	GetSession(ctx echo.Context, id string) error
	// Create a new global aspect.
	// (POST /sessions/{id}/aspects)
	CreateAspect(ctx echo.Context, id string) error
	// Delete an aspect.
	// (DELETE /sessions/{id}/aspects/{aspectId})
	DeleteAspect(ctx echo.Context, id string, aspectId string) error
	// Create a new character.
	// (POST /sessions/{id}/characters)
	CreateCharacter(ctx echo.Context, id string) error
	// Delete a character.
	// (DELETE /sessions/{id}/characters/{characterId})
	DeleteCharacter(ctx echo.Context, id string, characterId string) error
	// Create a new aspect bound to a specific character.
	// (POST /sessions/{id}/characters/{characterId}/aspects)
	CreateCharacterAspect(ctx echo.Context, id string, characterId string) error
	// Update Fate Points for the character
	// (PUT /sessions/{id}/characters/{characterId}/fatepoints)
	UpdateFatePoints(ctx echo.Context, id string, characterId string) error
	// Retrieve version information
	// (GET /version-info)
	GetVersionInfo(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateAuthToken converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAuthToken(ctx echo.Context) error {
	var err error

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateAuthToken(ctx)
	return err
}

// CreateSession converts echo context to params.
func (w *ServerInterfaceWrapper) CreateSession(ctx echo.Context) error {
	var err error

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateSession(ctx)
	return err
}

// GetSession converts echo context to params.
func (w *ServerInterfaceWrapper) GetSession(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSession(ctx, id)
	return err
}

// CreateAspect converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAspect(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateAspect(ctx, id)
	return err
}

// DeleteAspect converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteAspect(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "aspectId" -------------
	var aspectId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "aspectId", runtime.ParamLocationPath, ctx.Param("aspectId"), &aspectId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter aspectId: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteAspect(ctx, id, aspectId)
	return err
}

// CreateCharacter converts echo context to params.
func (w *ServerInterfaceWrapper) CreateCharacter(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateCharacter(ctx, id)
	return err
}

// DeleteCharacter converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteCharacter(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "characterId" -------------
	var characterId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "characterId", runtime.ParamLocationPath, ctx.Param("characterId"), &characterId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter characterId: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteCharacter(ctx, id, characterId)
	return err
}

// CreateCharacterAspect converts echo context to params.
func (w *ServerInterfaceWrapper) CreateCharacterAspect(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "characterId" -------------
	var characterId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "characterId", runtime.ParamLocationPath, ctx.Param("characterId"), &characterId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter characterId: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateCharacterAspect(ctx, id, characterId)
	return err
}

// UpdateFatePoints converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateFatePoints(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// ------------- Path parameter "characterId" -------------
	var characterId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "characterId", runtime.ParamLocationPath, ctx.Param("characterId"), &characterId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter characterId: %s", err))
	}

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateFatePoints(ctx, id, characterId)
	return err
}

// GetVersionInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetVersionInfo(ctx echo.Context) error {
	var err error

	ctx.Set(BearerScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetVersionInfo(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/auth/new", wrapper.CreateAuthToken)
	router.POST(baseURL+"/sessions", wrapper.CreateSession)
	router.GET(baseURL+"/sessions/:id", wrapper.GetSession)
	router.POST(baseURL+"/sessions/:id/aspects", wrapper.CreateAspect)
	router.DELETE(baseURL+"/sessions/:id/aspects/:aspectId", wrapper.DeleteAspect)
	router.POST(baseURL+"/sessions/:id/characters", wrapper.CreateCharacter)
	router.DELETE(baseURL+"/sessions/:id/characters/:characterId", wrapper.DeleteCharacter)
	router.POST(baseURL+"/sessions/:id/characters/:characterId/aspects", wrapper.CreateCharacterAspect)
	router.PUT(baseURL+"/sessions/:id/characters/:characterId/fatepoints", wrapper.UpdateFatePoints)
	router.GET(baseURL+"/version-info", wrapper.GetVersionInfo)

}
