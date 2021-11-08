package radio

import "log"

func NewState(uuid string) *State {
	boolDefault := false
	stringDefault := ""
	intDefault := 0
	return &State{
		IsMuted:  &boolDefault,
		Metadata: &stringDefault,
		Power:    &boolDefault,
		Preset:   -1,
		UUID:     uuid,
		Volume:   &intDefault,
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
	if ss.Preset != 0 {
		s.Preset = ss.Preset
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

func (h *Hub) AddClient(client *chan State) {
	h.stateOPS <- func(m map[*chan State]bool) {
		m[client] = true
		log.Println("Hub.stateLoop: client registered")
	}
}

func (h *Hub) RemoveClient(client *chan State) {
	h.stateOPS <- func(m map[*chan State]bool) {
		if _, ok := m[client]; ok {
			delete(m, client)
			close(*client)
			log.Println("Hub.stateLoop: client unregistered")
		}
	}
}

func (h *Hub) emitState(state *State) {
	h.stateOPS <- func(m map[*chan State]bool) {
		for client := range m {
			select {
			case *client <- *state:
			default:
				delete(m, client)
				close(*client)
				log.Println("Hub.stateLoop: client deleted")
			}
		}
	}
}

func (h *Hub) stateLoop() {
	log.Println("Hub.stateLoop: started")
	m := make(map[*chan State]bool)
	for op := range h.stateOPS {
		op(m)
	}
}
