package state

import "fmt"

const (
	StatusConnecting = "Connecting"
	StatusPlaying    = "Playing"
	StatusStopped    = "Stopped"
	StatusUnknown    = ""

	AudioSourceInternetRadio = "Internet radio"
)

const (
	ChangedAll = 1 << iota
	ChangedAudioSource
	ChangedIsMuted
	ChangedMetadata
	ChangedPower
	ChangedPresetNumber
	ChangedPresets
	ChangedStatus
	ChangedTitle
	ChangedTitleNew
	ChangedURL
	ChangedURLNew
	ChangedUUID
	ChangedVolume
)

type (
	Pub interface {
		Publish(state State, changed int)
		Subscribe(uuid string) (<-chan Message, func())
	}

	Message struct {
		State   State
		Changed int
	}

	Middleware interface {
		Apply(*Fragment)
	}

	MiddlewarePub interface {
		Publish()
		Subscribe() (<-chan struct{}, func())
	}

	Preset struct {
		Number   int    `json:"number"`    // Number is the preset number.
		Title    string `json:"title"`     // Title of the preset.
		TitleNew string `json:"title_new"` // TitleNew is the overridden title.
		URL      string `json:"url"`       // URL of the preset.
		URLNew   string `json:"url_new"`   // URLNew is the overridden URL.
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
		TitleNew     string   `json:"title_new"`     // TitleNew is the overridden title.
		URL          string   `json:"url"`           // URL of the stream that is currently selected.
		URLNew       string   `json:"url_new"`       // URLNew is the overridden URL.
		UUID         string   `json:"uuid"`          // UUID of the radio.
		Volume       int      `json:"volume"`        // Volume of the radio.
	}

	Partial struct {
		AudioSource  *string  `json:"audio_source,omitempty"`
		IsMuted      *bool    `json:"is_muted,omitempty"`
		Metadata     *string  `json:"metadata,omitempty"`
		Power        *bool    `json:"power,omitempty"`
		PresetNumber *int     `json:"preset_number,omitempty"`
		Presets      []Preset `json:"presets,omitempty"`
		Status       *Status  `json:"status,omitempty"`
		Title        *string  `json:"title,omitempty"`
		TitleNew     *string  `json:"title_new,omitempty"`
		URL          *string  `json:"url,omitempty"`
		URLNew       *string  `json:"url_new,omitempty"`
		UUID         string   `json:"uuid"`
		Volume       *int     `json:"volume,omitempty"`
	}

	Fragment struct {
		AudioSource *string
		IsMuted     *bool
		Metadata    *string
		Power       *bool
		Presets     []Preset
		Status      *Status
		Title       *string
		TitleNew    *string
		URL         *string
		URLNew      *string
		UUID        string
		Volume      *int
	}

	Status string
)

func NewPreset(number int, title, url string) Preset {
	return Preset{
		Number: number,
		Title:  title,
		URL:    url,
	}
}

func New(uuid, name, modelName, modelNumber string) State {
	return State{
		ModelName:   modelName,
		ModelNumber: modelNumber,
		Name:        name,
		Status:      StatusUnknown,
		UUID:        uuid,
	}
}

func NewPartial(uuid string) Partial {
	return Partial{
		UUID: uuid,
	}
}

func NewFragment(uuid string) Fragment {
	return Fragment{
		UUID: uuid,
	}
}

func IsChangedAll(changed int) bool {
	return changed&ChangedAll == ChangedAll
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
