package radio

import "errors"

var (
	ErrPresetInvalid = errors.New("invalid preset")
	ErrVolumeInvalid = errors.New("invalid volume")
	ErrRadioNotFound = errors.New("radio not found")
	ErrDiscovering   = errors.New("radios are already being discovered")
)
