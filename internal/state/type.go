package state

import (
	"fmt"
	"math"
)

const (
	StatusConnecting Status = "Connecting"
	StatusPlaying    Status = "Playing"
	StatusStopped    Status = "Stopped"
	StatusUnknown    Status = ""
)

const AudioSourceInternetRadio = "Internet radio"

const ChangedAll Changed = math.MaxInt

const (
	ChangedNone                = 0
	ChangedAudioSource Changed = 1 << iota
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
	Preset struct {
		Number   int    `json:"number" validate:"required"`    // Number is the preset number.
		Title    string `json:"title" validate:"required"`     // Title of the preset.
		TitleNew string `json:"title_new" validate:"required"` // TitleNew is the overridden title.
		URL      string `json:"url" validate:"required"`       // URL of the preset.
		URLNew   string `json:"url_new" validate:"required"`   // URLNew is the overridden URL.
	} //	@name	state.Preset

	State struct {
		AudioSource  string   `json:"audio_source" validate:"required"`  // AudioSource is the audio source.
		AudioSources []string `json:"audio_sources" validate:"required"` // AudioSources is the list of available audio sources.
		IsMuted      bool     `json:"is_muted" validate:"required"`      // IsMuted represents if the radio is muted.
		Metadata     string   `json:"metadata" validate:"required"`      // Metadata of the current playing stream.
		ModelName    string   `json:"model_name" validate:"required"`    // ModelName is the model name of the device.
		ModelNumber  string   `json:"model_number" validate:"required"`  // ModelNumber is the model number of the device.
		Name         string   `json:"name" validate:"required"`          // Name of the radio.
		Power        bool     `json:"power" validate:"required"`         // Power represents if the radio is not in standby.
		PresetNumber int      `json:"preset_number" validate:"required"` // PresetNumber is the current preset that is playing.
		Presets      []Preset `json:"presets" validate:"required"`       // Presets of the radio.
		Status       Status   `json:"status" validate:"required"`        // Status is either playing, connecting, or stopped.
		Title        string   `json:"title" validate:"required"`         // Title of the current playing stream.
		TitleNew     string   `json:"title_new" validate:"required"`     // TitleNew is the overridden title.
		URL          string   `json:"url" validate:"required"`           // URL of the stream that is currently selected.
		URLNew       string   `json:"url_new" validate:"required"`       // URLNew is the overridden URL.
		UUID         string   `json:"uuid" validate:"required"`          // UUID of the radio.
		Volume       int      `json:"volume" validate:"required"`        // Volume of the radio.
	} //	@name	state.State

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

	Status string

	Changed int
)

func NewPreset(number int, title, url string) Preset {
	return Preset{
		Number: number,
		Title:  title,
		URL:    url,
	}
}

func New(uuid, name, modelName, modelNumber string, audioSources []string) State {
	return State{
		AudioSources: audioSources,
		ModelName:    modelName,
		ModelNumber:  modelNumber,
		Name:         name,
		Status:       StatusUnknown,
		UUID:         uuid,
	}
}

func NewPartial(uuid string) Partial {
	return Partial{
		UUID: uuid,
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
	case status == string(StatusConnecting):
		return StatusConnecting
	case status == string(StatusPlaying):
		return StatusPlaying
	case status == string(StatusStopped):
		return StatusStopped
	default:
		return StatusUnknown
	}
}
