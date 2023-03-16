package background

import "context"

type Background interface {
	Background(ctx context.Context, doneC chan<- struct{})
}

func Run(ctx context.Context, backgrounds []Background) <-chan struct{} {
	doneFanIn := make(chan struct{}, len(backgrounds))
	running := 0

	// Start backgrounds
	for _, background := range backgrounds {
		go background.Background(ctx, doneFanIn)
		running++
	}

	done := make(chan struct{})

	go func() {
		// Wait for context
		<-ctx.Done()

		// Wait for backgrounds
		for i := 0; i < running; i++ {
			<-doneFanIn
		}

		close(done)
	}()

	return done
}

type BackgroundFunction func(ctx context.Context)

type Function struct {
	fn BackgroundFunction
}

func NewFunction(fn BackgroundFunction) Function {
	return Function{fn: fn}
}

func (bf Function) Background(ctx context.Context, doneC chan<- struct{}) {
	doneC <- struct{}{}
	bf.fn(ctx)
}
