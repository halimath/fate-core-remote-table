package ingress

import (
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createaspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createcharacteraspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/deleteaspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/joinsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/loadsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/updatefatepoints"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/fate-core-remote-table/backend/internal/ingress/rest"
	"github.com/halimath/fate-core-remote-table/backend/internal/ingress/web"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func Provide(cfg config.Config, logger kvlog.Logger, version, commit string,
	tokenHandler auth.TokenHandler,
	createSession createsession.Port,
	loadSession loadsession.Port,
	joinSession joinsession.Port,
	createAspect createaspect.Port,
	createCharacterAspect createcharacteraspect.Port,
	deleteAspect deleteaspect.Port,
	updateFatePoints updatefatepoints.Port,
) http.Handler {
	if cfg.DevMode {
		response.DevMode = true
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", rest.Provide(cfg, logger, version, commit, tokenHandler, createSession, loadSession, joinSession, createAspect, createCharacterAspect, deleteAspect, updateFatePoints))
	mux.Handle("/", web.Provide())

	return kvlog.Middleware(logger, true)(mux)
}
