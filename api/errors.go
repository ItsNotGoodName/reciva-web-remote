package api

import "errors"

var ErrPresetNewNameInvalid = errors.New("preset new name invalid")
var ErrPresetNotFound = errors.New("preset not found")
var ErrNameAlreadyExists = errors.New("name already exists")
