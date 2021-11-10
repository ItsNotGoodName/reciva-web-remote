package store

func NewStore(configFile string) (*Store, error) {
	m, err := readConfig(configFile)

	s := Store{presetOp: make(chan func(map[string]Preset)), configFile: configFile}
	if err != nil {
		s.readonly = true
	}

	go s.StoreLoop(m)

	return &s, err
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

func (s *Store) IsReadOnly() bool {
	return s.readonly
}
