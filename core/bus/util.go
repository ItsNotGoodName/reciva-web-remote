package bus

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
)

func writeCommand(ctx context.Context, writeC chan<- dto.Command, command dto.Command) {
	select {
	case <-ctx.Done():
		return
	case writeC <- command:
	}
}
