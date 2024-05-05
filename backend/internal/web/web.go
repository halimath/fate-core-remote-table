package web

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecase"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/httputils/requesturi"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

var (
	//go:embed public
	staticFiles embed.FS
)

func Provide(cfg config.Config, logger kvlog.Logger, version, commit string,
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

	if cfg.DevMode {
		response.DevMode = true
	}

	// TODO: Inject
	authProvider := auth.Provide(cfg)

	mux := http.NewServeMux()
	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", newAuthMux(authProvider)))
	mux.Handle("/api/sessions/", http.StripPrefix("/api/sessions", authMiddleware(authProvider)(newSessionAPIHandler(
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

func bindBody(r *http.Request, payload any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, payload)
}
