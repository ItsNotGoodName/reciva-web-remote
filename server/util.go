package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
)

// GetPresetURLS returns all urls for presets
func GetPresetURLS(p *api.PresetAPI) []string {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	presets, err := p.ReadPresets(ctx)
	cancel()
	if err != nil {
		log.Fatal("server.GetPresets(ERROR):", err)
	}
	urls := make([]string, len(presets))
	for i := range presets {
		urls[i] = presets[i].URL
	}
	return urls
}

// PrintAddresses prints listening addresses.
func PrintAddresses(port string) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("server.PrintAddresses(ERROR):", err)
		return
	}
	message := "\nNavigate to one of the following addresses:\n"
	for i := range addr {
		ip := net.ParseIP(strings.Split(addr[i].String(), "/")[0])
		if ip != nil && ip.To4() != nil {
			message = message + "\thttp://" + ip.String() + ":" + port + "\n"
		}
	}
	fmt.Println(message)
}
