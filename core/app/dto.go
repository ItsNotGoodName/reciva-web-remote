package app

import (
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type Command struct {
	Type Type        `json:"type"`
	Slug interface{} `json:"slug"`
}

type Type string

const (
	TypeError            = Type("error")
	TypeState            = Type("state")
	TypeStatePartial     = Type("state.partial")
	TypeStateSubscribe   = Type("state.subscribe")
	TypeStateUnsubscribe = Type("state.unsubscribe")
)

func NewErrorCommand(err error) Command {
	return Command{
		Type: TypeError,
		Slug: err.Error(),
	}
}

func NewStateCommand(state *state.State) Command {
	return Command{
		Type: TypeState,
		Slug: state,
	}
}

func NewStatePartialCommand(partial state.Partial) Command {
	return Command{
		Type: TypeStatePartial,
		Slug: partial,
	}
}
