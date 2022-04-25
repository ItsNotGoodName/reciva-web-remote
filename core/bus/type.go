package bus

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type (
	Type string

	Command struct {
		Type Type        `json:"type"`
		Slug interface{} `json:"slug"`
	}

	Service interface {
		Handle(ctx context.Context, readC <-chan Command, writeC chan<- Command)
	}
)

const (
	TypeError            = Type("error")
	TypeState            = Type("state")
	TypeStatePartial     = Type("state.partial")
	TypeStateSubscribe   = Type("state.subscribe")
	TypeStateUnsubscribe = Type("state.unsubscribe")
)

func newErrorCommand(err error) Command {
	return Command{
		Type: TypeError,
		Slug: err.Error(),
	}
}

func newStateCommand(state *state.State) Command {
	return Command{
		Type: TypeState,
		Slug: state,
	}
}

func newStatePartialCommand(partial state.Partial) Command {
	return Command{
		Type: TypeStatePartial,
		Slug: partial,
	}
}
