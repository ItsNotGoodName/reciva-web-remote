package radio

import (
	"context"
	"regexp"

	"github.com/huin/goupnp"
)

func normalizeVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}
	return volume
}

var uuidReg = regexp.MustCompile(`(?m)\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func getServiceClientUUID(c *goupnp.ServiceClient) (string, bool) {
	uuid := uuidReg.FindString(c.Location.String())
	if uuid == "" {
		return uuid, false
	}
	return uuid, true
}

func presetMutator(ctx context.Context, p *Preset) {
	p.Name = p.Title
}
