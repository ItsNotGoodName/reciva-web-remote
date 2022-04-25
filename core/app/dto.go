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
	TypeStateSubscribe   = Type("state.subscribe")
	TypeStateUnsubscribe = Type("state.unsubscribe")
)

func NewStateCommand(state *state.State) Command {
	return Command{
		Type: TypeState,
		Slug: state,
	}
}

func NewErrorCommand(err error) Command {
	return Command{
		Type: TypeError,
		Slug: err.Error(),
	}
}
