package radio

import "errors"

var (
	ErrInvalidPreset = errors.New("invalid preset")
	ErrInvalidVolume = errors.New("invalid volume")
	ErrRadioNotFound = errors.New("radio not found")
	ErrDiscovering   = errors.New("radios are already being discovered")
)
