package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WaitForShutdown blocks until the process receives SIGINT or SIGTERM
// Returns a context that is canceled when the shutdown signal is received
func WaitForShutdown(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	return ctx
}
