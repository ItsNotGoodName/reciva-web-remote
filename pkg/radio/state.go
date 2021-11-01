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

func (s *State) Merge(ss *State) {
	if ss.IsMuted != nil {
		s.IsMuted = ss.IsMuted
	}
	if ss.Metadata != nil {
		s.Metadata = ss.Metadata
	}
	if ss.Name != "" {
		s.Name = ss.Name
	}
	if ss.NumPresets != 0 {
		s.NumPresets = ss.NumPresets
	}
	if ss.Power != nil {
		s.Power = ss.Power
	}
	if len(ss.Presets) != 0 {
		s.Presets = ss.Presets
	}
	if ss.PresetNum != 0 {
		s.PresetNum = ss.PresetNum
	}
	if ss.State != "" {
		s.State = ss.State
	}
	if ss.Title != "" {
		s.Title = ss.Title
	}
	if ss.URL != "" {
		s.URL = ss.URL
	}
	if ss.Volume != nil {
		s.Volume = ss.Volume
	}
}
