package app

import "context"

func writeCommand(ctx context.Context, writeC chan<- Command, command Command) {
	select {
	case <-ctx.Done():
		return
	case writeC <- command:
	}
}
