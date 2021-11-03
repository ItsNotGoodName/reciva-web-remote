package api

import "errors"

var ErrStreamNotFound = errors.New("stream not found")
var ErrNameAlreadyExists = errors.New("name already exists")
