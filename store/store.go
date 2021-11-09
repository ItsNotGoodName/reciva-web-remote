package store

func NewStore(configFile string) (*Store, error) {
	m, err := readConfig(configFile)

	if err != nil {
		return nil, err
	}

	s := Store{presetOp: make(chan func(map[string]Preset)), configFile: configFile}
	go s.StoreLoop(m)

	return &s, nil
}

func (s *Store) StoreLoop(presets map[string]Preset) {
	for {
		select {
		case f := <-s.presetOp:
			f(presets)
		case f := <-s.configOp:
			presets = f(presets)
		}
	}
}
