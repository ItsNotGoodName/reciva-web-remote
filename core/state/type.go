package state

import "fmt"

const (
	StatusConnecting         = "Connecting"
	StatusPlaying            = "Playing"
	StatusStopped            = "Stopped"
	StatusUnknown            = ""
	AudioSourceInternetRadio = "Internet radio"
)

type (
	StatePub interface {
		Publish(State)
		Subscribe(uuid string) (<-chan State, func())
	}

	Middleware interface {
		Apply(*Fragment)
	}

	MiddlewarePub interface {
		Publish()
		Subscribe() (<-chan struct{}, func())
	}

	State struct {
		AudioSource  string   `json:"audio_source"`  // AudioSource is the audio source.
		AudioSources []string `json:"audio_sources"` // AudioSources is the list of available audio sources.
		IsMuted      bool     `json:"is_muted"`      // IsMuted represents if the radio is muted.
		Metadata     string   `json:"metadata"`      // Metadata of the current playing stream.
		ModelName    string   `json:"model_name"`    // ModelName is the model name of the device.
		ModelNumber  string   `json:"model_number"`  // ModelNumber is the model number of the device.
		Name         string   `json:"name"`          // Name of the radio.
		NewTitle     string   `json:"new_title"`     // NewTitle is the overridden title.
		NewURL       string   `json:"new_url"`       // NewURL is the overridden URL.
		Power        bool     `json:"power"`         // Power represents if the radio is not in standby.
		PresetNumber int      `json:"preset_number"` // PresetNumber is the current preset that is playing.
		Presets      []Preset `json:"presets"`       // Presets of the radio.
		Status       Status   `json:"status"`        // Status is either playing, connecting, or stopped.
		Title        string   `json:"title"`         // Title of the current playing stream.
		URL          string   `json:"url"`           // URL of the stream that is currently selected.
		UUID         string   `json:"uuid"`          // UUID of the radio.
		Volume       int      `json:"volume"`        // Volume of the radio.
	}

	Fragment struct {
		AudioSource *string
		IsMuted     *bool
		Metadata    *string
		NewTitle    *string
		NewURL      *string
		Power       *bool
		Presets     []Preset
		Status      *Status
		Title       *string
		URL         *string
		UUID        string
		Volume      *int
	}

	Preset struct {
		Number   int    `json:"number"`    // Number is the preset number.
		Title    string `json:"title"`     // Title of the preset.
		NewTitle string `json:"new_title"` // NewTitle is the overridden title.
		URL      string `json:"url"`       // URL of the preset.
	}

	Status string
)

func New(uuid, name, modelName, modelNumber string) State {
	return State{
		ModelName:   modelName,
		ModelNumber: modelNumber,
		Name:        name,
		Status:      StatusUnknown,
		UUID:        uuid,
	}
}

func NewPreset(number int, title, url string) Preset {
	return Preset{
		Number: number,
		Title:  title,
		URL:    url,
	}
}

func ValidPresetNumber(s *State, preset int) error {
	if preset < 1 || preset > len(s.Presets) {
		return fmt.Errorf("invalid preset number: %d", preset)
	}

	return nil
}

func ValidAudioSource(s *State, audioSource string) error {
	for _, source := range s.AudioSources {
		if source == audioSource {
			return nil
		}
	}

	return fmt.Errorf("invalid audio source: %s", audioSource)
}

func NormalizeVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}

	return volume
}

func ParsePresetsCount(presetsCount int) (int, error) {
	presetsCount = presetsCount - 2
	if presetsCount < 1 {
		return 0, fmt.Errorf("invalid presets count: %d", presetsCount)
	}

	return presetsCount, nil
}

func ParseStatus(status string) Status {
	switch {
	case status == StatusConnecting:
		return StatusConnecting
	case status == StatusPlaying:
		return StatusPlaying
	case status == StatusStopped:
		return StatusStopped
	default:
		return StatusUnknown
	}
}
