package dto

import (
	"context"
)

type (
	CommandType string

	Command struct {
		Type CommandType `json:"type"`
		Slug interface{} `json:"slug"`
	}

	Bus interface {
		Handle(ctx context.Context, readC <-chan Command, writeC chan<- Command)
	}
)

const (
	CommandTypeError            = CommandType("error")
	CommandTypeState            = CommandType("state")
	CommandTypeStatePartial     = CommandType("state.partial")
	CommandTypeStateSubscribe   = CommandType("state.subscribe")
	CommandTypeStateUnsubscribe = CommandType("state.unsubscribe")
)
