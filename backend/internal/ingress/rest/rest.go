package rest

import (
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecase"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func Provide(
	cfg config.Config, logger kvlog.Logger, version, commit string,
	tokenHandler auth.TokenHandler,
	createSession usecase.CreateSession,
	loadSession usecase.LoadSession,
	joinSession usecase.JoinSession,
	createAspect usecase.CreateAspect,
	createCharacterAspect usecase.CreateCharacterAspect,
	deleteAspect usecase.DeleteAspect,
	updateFatePoints usecase.UpdateFatePoints,
) http.Handler {
	versionInfo := VersionInfo{
		Version:    version,
		Commit:     commit,
		ApiVersion: "1.0.0",
	}

	mux := http.NewServeMux()
	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", newAuthMux(tokenHandler)))
	mux.Handle("/api/sessions/", http.StripPrefix("/api/sessions", authMiddleware(tokenHandler)(newSessionAPIHandler(
		cfg,
		createSession,
		loadSession,
		joinSession,
		createAspect,
		createCharacterAspect,
		deleteAspect,
		updateFatePoints,
	))))
	mux.HandleFunc("GET /api/version-info", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, r, versionInfo)
	})

	return mux
}
