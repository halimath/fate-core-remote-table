package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/halimath/fate-core-remote-table/backend/internal/boundary"
	"github.com/halimath/fate-core-remote-table/backend/internal/control"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/kvlog"

	_ "time/tzdata"
)

var (
	Version string = "0.0.0"
	Commit  string = "local"
)

func main() {
	ctx := context.Background()

	cfg := config.Provide(ctx)
	if cfg.DevMode {
		kvlog.L = kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.ConsoleFormatter()))
	} else {
		kvlog.L = kvlog.New(kvlog.NewSyncHandler(os.Stdout, kvlog.JSONLFormatter()))
	}
	kvlog.L.AddHook(kvlog.TimeHook)

	controller := control.Provide(cfg)
	mux := boundary.Provide(cfg, controller, kvlog.L, Version, Commit)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
	}

	kvlog.L.Logs("startup", kvlog.WithKV("version", Version), kvlog.WithKV("commit", Commit))

	termChan := make(chan int, 1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-signalCh

		kvlog.L.Logs("receivedSignal", kvlog.WithKV("signal", s))
		httpServer.Close()

		termChan <- 0
	}()

	go func() {
		kvlog.L.Logs("httpListen", kvlog.WithKV("addr", ":8080"))
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			kvlog.L.Logs("httpServerFailedToStart", kvlog.WithErr(err))
			termChan <- 1
		}
	}()

	exitCode := <-termChan
	kvlog.L.Logs("exit", kvlog.WithKV("code", exitCode))
	os.Exit(exitCode)
}
