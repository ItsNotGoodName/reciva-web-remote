package radio

func NewState(uuid string) *State {
	boolDefault := false
	stringDefault := ""
	intDefault := 0
	return &State{
		IsMuted:  &boolDefault,
		Metadata: &stringDefault,
		Power:    &boolDefault,
		UUID:     uuid,
		Volume:   &intDefault,
	}
}
