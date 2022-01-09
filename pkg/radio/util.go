package radio

import (
	"fmt"
	"regexp"

	"github.com/huin/goupnp"
)

var uuidReg = regexp.MustCompile(`(?m)\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func getServiceClientUUID(c goupnp.ServiceClient) (string, error) {
	uuid := uuidReg.FindString(c.Location.String())
	if uuid == "" {
		return "", fmt.Errorf("could not find UUID in location: %s", c.Location)
	}

	return uuid, nil
}

func normalizeVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}
	return volume
}
