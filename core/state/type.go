package state

const (
	StatusConnecting = "Connecting"
	StatusPlaying    = "Playing"
	StatusStopped    = "Stopped"
	StatusUnknown    = ""
)

type (
	FragmentPub interface {
		Publish(Fragment)
	}

	State struct {
		AudioSource  string   `json:"audio_source"`  // AudioSource is the audio source.
		AudioSources []string `json:"audio_sources"` // AudioSources is the list of available audio sources.
		IsMuted      bool     `json:"is_muted"`      // IsMuted represents if the radio is muted.
		Metadata     string   `json:"metadata"`      // Metadata of the current playing stream.
		ModelName    string   `json:"model_name"`    // ModelName is the model name of the device.
		ModelNumber  string   `json:"model_number"`  // ModelNumber is the model number of the device.
		Name         string   `json:"name"`          // Name of the radio.
		Power        bool     `json:"power"`         // Power represents if the radio is not in standby.
		PresetNumber int      `json:"preset_number"` // PresetNumber is the current preset that is playing.
		Presets      []Preset `json:"presets"`       // Presets of the radio.
		Status       Status   `json:"status"`        // Status is either playing, connecting, or stopped.
		Title        string   `json:"title"`         // Title of the current playing stream.
		URL          string   `json:"url"`           // URL of the stream that is currently selected.
		UUID         string   `json:"uuid"`          // UUID of the radio.
		Volume       int      `json:"volume"`        // Volume of the radio.
	}

	// Fragment is the fragment of the state that changes.
	Fragment struct {
		AudioSource  *string  `json:"audio_source,omitempty"`
		IsMuted      *bool    `json:"is_muted,omitempty"`
		Metadata     *string  `json:"metadata,omitempty"`
		Power        *bool    `json:"power,omitempty"`
		PresetNumber *int     `json:"preset_number,omitempty"`
		Presets      []Preset `json:"presets,omitempty"`
		Status       *Status  `json:"status,omitempty"`
		Title        *string  `json:"title,omitempty"`
		URL          *string  `json:"url,omitempty"`
		UUID         string   `json:"uuid"`
		Volume       *int     `json:"volume,omitempty"`
	}

	Preset struct {
		Number int    `json:"number"` // Number is the preset number.
		Title  string `json:"title"`  // Title of the preset.
		URL    string `json:"url"`    // URL of the preset.
	}

	Status string
)

func NewFragment(uuid string) Fragment {
	return Fragment{
		UUID: uuid,
	}
}

func NewPreset(number int, title, url string) Preset {
	return Preset{
		Number: number,
		Title:  title,
		URL:    url,
	}
}
