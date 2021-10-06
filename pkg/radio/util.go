package radio

import (
	"regexp"

	"github.com/huin/goupnp"
)

var uuidReg = regexp.MustCompile(`(?m)\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func getServiceClientUUID(c *goupnp.ServiceClient) (string, bool) {
	uuid := uuidReg.FindString(c.Location.String())
	if uuid == "" {
		return uuid, false
	}
	return uuid, true
}

func IsValidVolume(volume int) bool {
	return volume >= 0 && volume <= 100
}
