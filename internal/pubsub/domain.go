package pubsub

import (
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
)

// State

const TopicState Topic = "state"

type MessageState struct {
	State   state.State
	Changed state.Changed
}

func PublishState(p Pub, msg MessageState) {
	p.publish(TopicState, msg)
}

func ParseState(msg *Message) (MessageState, bool) {
	if msg.Topic == TopicState {
		return msg.Data.(MessageState), true
	}

	return MessageState{}, false
}

// Stale state hook

const TopicStaleStateHook Topic = "stale.state.hook"

type StaleStateHookMessage struct {
	Changed state.Changed
}

func PublishStaleStateHook(p Pub, data StaleStateHookMessage) {
	p.publish(TopicStaleStateHook, data)
}

func ParseStaleStateHook(msg *Message) (StaleStateHookMessage, bool) {
	if msg.Topic == TopicStaleStateHook {
		return msg.Data.(StaleStateHookMessage), true
	}
	return StaleStateHookMessage{}, false
}

// Discover

const TopicDiscover Topic = "discover"

func PublishDiscover(p Pub, discovering bool) {
	p.publish(TopicDiscover, discovering)
}

func ParseDiscover(msg *Message) (bool, bool) {
	if msg.Topic == TopicDiscover {
		return msg.Data.(bool), true
	}
	return false, false
}

// Stale radios

const TopicStaleRadios Topic = "stale.radios"

func PublishStaleRadios(p Pub) {
	p.publish(TopicStaleRadios, nil)
}

func ParseStaleRadios(msg *Message) bool {
	return msg.Topic == TopicStaleRadios
}

// Error

const TopicError Topic = "error"

func PublishError(p Pub, err error) {
	p.publish(TopicError, err)
}

func ParseError(msg *Message) (error, bool) {
	if msg.Topic == TopicError {
		return msg.Data.(error), true
	}
	return nil, false
}

// Stale

const TopicStale Topic = "stale"

type Stale = string

const StaleRadios Stale = "radios"
const StalePresets Stale = "presets"
