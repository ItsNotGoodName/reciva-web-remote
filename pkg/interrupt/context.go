package interrupt

import (
	"context"
	"os"
	"os/signal"
)

// Context will return a context that will be cancelled on os.interrupt signal.
func Context() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
