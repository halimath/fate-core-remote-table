package boundary

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/boundary/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/control"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

var (
	//go:embed public
	staticFiles embed.FS
)

func Provide(cfg config.Config, ctrl control.SessionController, logger kvlog.Logger, version, commit string) http.Handler {
	authProvider := auth.Provide(cfg)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", HandlerWithOptions(&restHandler{
		versionInfo: VersionInfo{
			Version:    version,
			Commit:     commit,
			ApiVersion: "1.0.0",
		},
		controller:   ctrl,
		authProvider: authProvider,
	}, StdHTTPServerOptions{
		Middlewares:      []MiddlewareFunc{auth.Middleware(authProvider)},
		ErrorHandlerFunc: handleError,
	})))

	staticFilesFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		panic(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticFilesFS)))

	return kvlog.Middleware(logger, true)(mux)

	// e.Pre(middleware.Rewrite(map[string]string{
	// 	"/join/*":    "/",
	// 	"/session/*": "/",
	// }))

}

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
	} else if errors.Is(err, control.ErrNotFound) {
		e.Code = http.StatusNotFound
	}

	response.JSON(w, r, e, response.StatusCode(e.Code))
}
