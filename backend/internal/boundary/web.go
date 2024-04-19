package boundary

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/boundary/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/control"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/httputils/requesturi"
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
		Middlewares: []MiddlewareFunc{auth.Middleware(authProvider)},
	})))

	staticFilesFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		panic(err)
	}

	pathRewriter, err := requesturi.RewritePath(map[string]string{
		"/join/*":    "/",
		"/session/*": "/",
	})
	if err != nil {
		panic(err)
	}

	mux.Handle("/", requesturi.Middleware(http.FileServer(http.FS(staticFilesFS)), pathRewriter))

	return kvlog.Middleware(logger, true)(mux)

	// e.Pre(middleware.Rewrite(map[string]string{
	// 	"/join/*":    "/",
	// 	"/session/*": "/",
	// }))

}
