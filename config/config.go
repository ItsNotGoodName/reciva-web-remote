package config

import (
	"flag"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

type Config struct {
	ConfigPath    string
	Port          int
	CPort         int
	EnablePresets bool
	Presets       []string
	APIURI        string
}

const DefaultPort = 8080
const DefaultFile = "reciva-web-remote.json"

func NewConfig() *Config {
	port := flag.Int("port", DefaultPort, "Listen port for web server.")
	cport := flag.Int("cport", goupnpsub.DefaultPort, "Listen port for UPnP notify server.")
	presetsFlag := flag.String("presets", "", "List of presets to host seperated by comma (ex. /01.m3u,/02.m3u).")
	config := flag.String("config", DefaultFile, "Path to config location.")
	enablePresets := false
	var presets []string

	flag.Parse()

	// Enable presets based on presetsFlag
	if *presetsFlag == "" {
		presets = make([]string, 0)
	} else {
		presets = strings.Split(*presetsFlag, ",")
		enablePresets = true
	}

	return &Config{
		Port:          *port,
		CPort:         *cport,
		EnablePresets: enablePresets,
		Presets:       presets,
		ConfigPath:    *config,
		APIURI:        "/v1",
	}
}
