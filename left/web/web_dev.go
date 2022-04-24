//go:build dev

package web

import "io/fs"

func FS() fs.FS {
	return empty{}
}

type empty struct{}

func (empty) Open(string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
