package config

import (
	"flag"
	"strings"
)

type Config struct {
	ConfigPath    string   `json:"configPath"`
	Port          int      `json:"port"`
	CPort         int      `json:"cport"`
	EnablePresets bool     `json:"enablePresets"`
	Presets       []string `json:"presets"`
}

func NewConfig() *Config {
	port := flag.Int("port", 8080, "Listen port for web server.")
	cport := flag.Int("cport", 8058, "Listen port for UPnP notify server.")
	presetsFlag := flag.String("presets", "", "List of presets to host seperated by comma (ex. /01.m3u,/02.m3u).")
	config := flag.String("config", "settings.json", "Path to config location.")
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
		ConfigPath:    *config,
	}
}
