package store

import "errors"

var ErrNotFound = errors.New("not found")
var ErrEmptyPresets = errors.New("empty presets")
var ErrReadOnly = errors.New("store started as read only")
