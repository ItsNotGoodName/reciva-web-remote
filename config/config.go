package config

import (
	"flag"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

type Config struct {
	APIURI         string
	CPort          int
	CPortFlag      bool
	ConfigPath     string
	Port           int
	PortFlag       bool
	Presets        []string
	PresetsEnabled bool
}

const (
	APIURI      = "/v1"
	DefaultFile = "reciva-web-remote.json"
	DefaultPort = 8080
)

func NewConfig() *Config {
	cportFlag := false
	portFlag := false
	PresetsEnabled := false
	var presets []string

	cport := flag.Int("cport", goupnpsub.DefaultPort, "Listen port for UPnP notify server.")
	port := flag.Int("port", DefaultPort, "Listen port for web server.")
	config := flag.String("config", DefaultFile, "Path to config location.")
	presetsFlag := flag.String("presets", "", "List of presets to host seperated by comma (ex. /01.m3u,/02.m3u).")

	flag.Parse()

	// Check if flag was specified for port and cport
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			portFlag = true
		}
		if f.Name == "cport" {
			cportFlag = true
		}
	})

	// Enable presets based on presetsFlag
	if *presetsFlag == "" {
		presets = make([]string, 0)
	} else {
		presets = strings.Split(*presetsFlag, ",")
	}

	return &Config{
		APIURI:         APIURI,
		CPort:          *cport,
		CPortFlag:      cportFlag,
		ConfigPath:     *config,
		Port:           *port,
		PortFlag:       portFlag,
		Presets:        presets,
		PresetsEnabled: PresetsEnabled,
	}
}
