package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/halimath/fate-core-remote-table/backend/internal"
	"github.com/halimath/kvlog"

	_ "time/tzdata"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-signalCh
		kvlog.L.Logs("received signal; initiating shut down", kvlog.WithKV("signal", s))
		cancel()
	}()

	os.Exit(internal.RunService(ctx))
}
