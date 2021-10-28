package radio

func NewState(uuid string) *State {
	boolDefault := false
	stringDefault := ""
	intDefault := 0
	return &State{
		UUID:     uuid,
		Power:    &boolDefault,
		IsMuted:  &boolDefault,
		Volume:   &intDefault,
		Metadata: &stringDefault,
	}
}
