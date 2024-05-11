package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createaspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createcharacteraspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/createsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/deleteaspect"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/joinsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/loadsession"
	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecases/updatefatepoints"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/fate-core-remote-table/backend/internal/ingress"
	"github.com/halimath/fate-core-remote-table/backend/internal/persistence"
	"github.com/halimath/kvlog"
)

var (
	Version string = "0.0.0"
	Commit  string = "local"
)

func RunService(ctx context.Context) int {
	cfg := config.Provide(ctx)
	if cfg.DevMode {
		kvlog.L = kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.ConsoleFormatter()))
	} else {
		kvlog.L = kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter()))
	}
	kvlog.L.AddHook(kvlog.TimeHook)

	tokenHandler := auth.Provide(cfg)

	sessionRepo := persistence.NewSessionRepository(cfg)
	createSession := createsession.Provide(sessionRepo)
	loadSession := loadsession.Provide(sessionRepo)
	joinSession := joinsession.Provide(sessionRepo)
	createAspect := createaspect.Provide(sessionRepo)
	createCharacterAspect := createcharacteraspect.Provide(sessionRepo)
	deleteAspect := deleteaspect.Provide(sessionRepo)
	updateFatePoints := updatefatepoints.Provide(sessionRepo)

	mux := ingress.Provide(cfg, kvlog.L, Version, Commit, tokenHandler, createSession,
		loadSession, joinSession, createAspect, createCharacterAspect,
		deleteAspect, updateFatePoints)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
	}

	kvlog.L.Logs("startup", kvlog.WithKV("version", Version), kvlog.WithKV("commit", Commit))

	termChan := make(chan int, 1)

	go func() {
		kvlog.L.Logs("http listen", kvlog.WithKV("addr", ":8080"))
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			kvlog.L.Logs("http server failed to start", kvlog.WithErr(err))
			termChan <- 1
			return
		}

		termChan <- 0
	}()

	go func() {
		<-ctx.Done()
		kvlog.L.Logs("context done; shutting done")
		httpServer.Close()
	}()

	exitCode := <-termChan
	close(termChan)

	kvlog.L.Logs("exit", kvlog.WithKV("code", exitCode))
	return exitCode
}
