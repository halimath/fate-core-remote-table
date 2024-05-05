package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/halimath/fate-core-remote-table/backend/internal/domain/usecase"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/fate-core-remote-table/backend/internal/repository"
	"github.com/halimath/fate-core-remote-table/backend/internal/web"
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

	sessionRepo := repository.NewSessionRepository(cfg)
	createSession := usecase.ProvideCreateSession(sessionRepo)
	loadSession := usecase.ProvideLoadSession(sessionRepo)
	joinSession := usecase.ProvideJoinSession(sessionRepo)
	createAspect := usecase.ProvideCreateAspect(sessionRepo)
	createCharacterAspect := usecase.ProvideCreateCharacterAspect(sessionRepo)
	deleteAspect := usecase.ProvideDeleteAspect(sessionRepo)
	updateFatePoints := usecase.ProvideUpdateFatePoints(sessionRepo)

	mux := web.Provide(cfg, kvlog.L, Version, Commit, createSession,
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
