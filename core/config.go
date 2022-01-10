package core

import "fmt"

var (
	ErrConfigReadonly = fmt.Errorf("config is readonly")
)

type Config struct {
	Presets []Preset `json:"presets"`
}
