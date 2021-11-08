package api

import "errors"

var ErrStreamNotFound = errors.New("stream not found")
var ErrStreamNameInvalid = errors.New("stream name invalid")
var ErrStreamContentInvalid = errors.New("stream content invalid")
var ErrPresetNotFound = errors.New("preset not found")
var ErrNameAlreadyExists = errors.New("name already exists")
