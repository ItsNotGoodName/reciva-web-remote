package core

import "fmt"

var (
	ErrPresetNotFound = fmt.Errorf("preset not found")
)

var (
	ErrHubDiscovering   = fmt.Errorf("hub is discovering")
	ErrHubServiceClosed = fmt.Errorf("hub service closed")
)

var (
	ErrRadioClosed   = fmt.Errorf("radio closed")
	ErrRadioNotFound = fmt.Errorf("radio not found")
)
