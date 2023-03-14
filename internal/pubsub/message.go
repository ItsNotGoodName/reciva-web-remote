package pubsub

import (
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
)

const ErrorTopic Topic = "error"

type ErrorMessage struct {
	Error error
}

const StateTopic Topic = "state"

type StateMessage struct {
	State   state.State
	Changed state.Changed
}

const StateHookStaleTopic Topic = "state.hook.stale"

type StateHookStaleMessage struct {
	Changed state.Changed
}

const DiscoverTopic Topic = "discover"

type DiscoverMessage struct {
	Discovering bool
}

const StaleTopic Topic = "stale"

type StaleMessage = model.Stale
