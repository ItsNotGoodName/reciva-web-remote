package radio

func NewState(uuid string) *State {
	boolDefault := false
	stringDefault := ""
	intDefault := 0
	return &State{
		IsMuted:   &boolDefault,
		Metadata:  &stringDefault,
		Power:     &boolDefault,
		PresetNum: -1,
		UUID:      uuid,
		Volume:    &intDefault,
	}
}
