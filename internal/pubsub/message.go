package pubsub

import "github.com/ItsNotGoodName/reciva-web-remote/internal/state"

const ErrorTopic Topic = "error"

type ErrorMessage struct {
	Error error
}

const StateTopic Topic = "state"

type StateMessage struct {
	State   state.State
	Changed state.Changed
}

const ForceStateChangedTopic Topic = "force.state.changed"
const DiscoverTopic Topic = "discover"

type DiscoverMessage struct {
	Discovering bool
}
