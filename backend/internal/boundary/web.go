package boundary

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/session"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/httputils/requesturi"
	"github.com/halimath/kvlog"
)

var (
	//go:embed public
	staticFiles embed.FS
)

func Provide(cfg config.Config, srv session.Service, logger kvlog.Logger, version, commit string) http.Handler {
	authProvider := auth.Provide(cfg)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", HandlerWithOptions(&restHandler{
		versionInfo: VersionInfo{
			Version:    version,
			Commit:     commit,
			ApiVersion: "1.0.0",
		},
		service:      srv,
		authProvider: authProvider,
	}, StdHTTPServerOptions{
		Middlewares: []MiddlewareFunc{authMiddleware(authProvider)},
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
}

func extractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", false
	}

	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return "", false
	}

	return authHeader[len("bearer "):], true

}

func authMiddleware(m *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractBearerToken(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			sub, err := m.Authorize(tokenString)
			if err != nil {
				kvlog.L.Logs("invalidAuthToken", kvlog.WithErr(err))
				next.ServeHTTP(w, r)
				return
			}

			r = r.WithContext(auth.SetUserID(r.Context(), sub))

			next.ServeHTTP(w, r)
		})
	}
}
