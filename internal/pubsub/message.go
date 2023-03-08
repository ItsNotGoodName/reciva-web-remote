package pubsub

import "github.com/ItsNotGoodName/reciva-web-remote/internal/state"

const StateTopic = "state"

type StateMessage struct {
	State   state.State
	Changed state.Changed
}

const PresetsMutatedTopic = "presets.mutated"
