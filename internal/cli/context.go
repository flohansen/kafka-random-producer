package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()

		<-sig
		os.Exit(1)
	}()

	return ctx
}
