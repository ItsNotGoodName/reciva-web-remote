package internal

import "fmt"

var (
	ErrPresetNotFound = fmt.Errorf("preset not found")
)

var ErrDiscovering = fmt.Errorf("discovery in progress")

var ErrHubClosed = fmt.Errorf("hub closed")

var (
	ErrRadioClosed   = fmt.Errorf("radio closed")
	ErrRadioNotFound = fmt.Errorf("radio not found")
)
