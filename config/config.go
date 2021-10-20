package config

import (
	"flag"
	"strings"
)

type Config struct {
	Port          int
	CPort         int
	EnablePresets bool
	Presets       []string
}

func NewConfig() *Config {
	port := flag.Int("port", 8080, "Listen port for web server.")
	cport := flag.Int("cport", 8058, "Listen port for UPnP notify server.")
	presetsFlag := flag.String("presets", "", "List of presets to host seperated by comma (ex. /01.m3u,/02.m3u).")
	enablePresets := false
	var presets []string

	flag.Parse()

	if *presetsFlag != "" {
		presets = strings.Split(*presetsFlag, ",")
		enablePresets = true
	}

	return &Config{
		Port:          *port,
		CPort:         *cport,
		EnablePresets: enablePresets,
		Presets:       presets,
	}
}
