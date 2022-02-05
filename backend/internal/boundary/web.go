package boundary

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/boundary/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/control"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/kvlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	//go:embed public
	staticFiles embed.FS
)

func Provide(cfg config.Config, ctrl control.SessionController, version, commit string) *echo.Echo {
	authProvider := auth.Provide(cfg)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Pre(middleware.Rewrite(map[string]string{
		"/join/*":    "/",
		"/session/*": "/",
	}))

	e.HTTPErrorHandler = handleError
	e.Use(loggingMiddleware)

	rest := e.Group("/api")
	rest.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	rest.Use(auth.Middleware(authProvider))
	handler := restHandler{
		versionInfo: VersionInfo{
			Version:    version,
			Commit:     commit,
			ApiVersion: "1.0.0",
		},
		controller:   ctrl,
		authProvider: authProvider,
	}
	RegisterHandlers(rest, &handler)

	staticFilesFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		panic(err)
	}

	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(staticFilesFS))))

	return e
}

func handleError(err error, ctx echo.Context) {
	e := Error{
		Error: err.Error(),
		Code:  http.StatusInternalServerError,
	}

	if httpError, ok := err.(*echo.HTTPError); ok {
		e.Code = httpError.Code
		e.Error = fmt.Sprintf("%s", httpError.Message)
	} else if errors.Is(err, control.ErrNotFound) {
		e.Code = http.StatusNotFound
	}

	ctx.JSON(e.Code, e)
}

func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		reqTime := time.Since(start)

		if err != nil {
			kvlog.Error(kvlog.Evt("requestError"), kvlog.Err(err))
			c.Error(err)
		}

		kvlog.Info(
			kvlog.Evt("request"),
			kvlog.KV("uri", c.Request().RequestURI),
			kvlog.KV("method", c.Request().Method),
			kvlog.KV("status", c.Response().Status),
			kvlog.Dur(reqTime),
		)

		return nil
	}
}
