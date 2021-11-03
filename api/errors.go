package api

import "errors"

var ErrStreamNotFound = errors.New("stream not found")
var ErrPresetNotFound = errors.New("preset not found")
var ErrNameAlreadyExists = errors.New("name already exists")
