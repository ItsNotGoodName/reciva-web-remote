package api

import "errors"

var ErrPresetNotFound = errors.New("preset not found")
var ErrNameAlreadyExists = errors.New("name already exists")
