package radio

import "errors"

var (
	ErrInvalidPreset = errors.New("invalid preset")
	ErrInvalidVolume = errors.New("invalid volume")
)
