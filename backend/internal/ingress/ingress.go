package ingress

import (
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecase"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/fate-core-remote-table/backend/internal/ingress/rest"
	"github.com/halimath/fate-core-remote-table/backend/internal/ingress/web"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func Provide(cfg config.Config, logger kvlog.Logger, version, commit string,
	tokenHandler auth.TokenHandler,
	createSession usecase.CreateSession,
	loadSession usecase.LoadSession,
	joinSession usecase.JoinSession,
	createAspect usecase.CreateAspect,
	createCharacterAspect usecase.CreateCharacterAspect,
	deleteAspect usecase.DeleteAspect,
	updateFatePoints usecase.UpdateFatePoints,
) http.Handler {
	if cfg.DevMode {
		response.DevMode = true
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", rest.Provide(cfg, logger, version, commit, tokenHandler, createSession, loadSession, joinSession, createAspect, createCharacterAspect, deleteAspect, updateFatePoints))
	mux.Handle("/", web.Provide())

	return kvlog.Middleware(logger, true)(mux)
}
