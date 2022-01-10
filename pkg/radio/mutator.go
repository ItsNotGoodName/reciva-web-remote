package radio

import "context"

type MutatorPort interface {
	GetTrigger() <-chan struct{}
	// Called when radio starts and when Change() triggers.
	Mutate(ctx context.Context, state *State) *State
	// MutateNewURL is called when the URL state changes.
	MutateNewURL(ctx context.Context, state *State) *State
}

type Mutator struct{}

func NewMutator() *Mutator {
	return &Mutator{}
}

func (m *Mutator) Mutate(ctx context.Context, state *State) *State {
	for i := range state.Presets {
		state.Presets[i].Name = state.Presets[i].Title
	}

	return &State{Presets: state.Presets}
}

func (m *Mutator) MutateNewURL(ctx context.Context, state *State) *State {
	return &State{}
}

func (m *Mutator) GetTrigger() <-chan struct{} {
	return make(<-chan struct{})
}
