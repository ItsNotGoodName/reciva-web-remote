package background

import "context"

type Background interface {
	Background(ctx context.Context, doneC chan<- struct{})
}

func Run(ctx context.Context, backgrounds []Background) {
	done := make(chan struct{}, len(backgrounds))
	running := 0

	// Start backgrounds
	for _, background := range backgrounds {
		go background.Background(ctx, done)
		running++
	}

	// Wait for backgrounds
	for i := 0; i < running; i++ {
		<-done
	}
}
