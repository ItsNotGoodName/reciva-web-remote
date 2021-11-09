package config

import (
	"flag"

	"github.com/ItsNotGoodName/go-upnpsub"
)

type Config struct {
	APIURI     string
	CPort      int
	ConfigFile string
	Port       int
}

const (
	APIURI      = "/v1"
	ConfigFile  = "reciva-web-remote.json"
	DefaultPort = 8080
)

func NewConfig(options ...func(*Config)) *Config {
	c := &Config{
		APIURI:     APIURI,
		CPort:      upnpsub.DefaultPort,
		ConfigFile: ConfigFile,
		Port:       DefaultPort,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithFlag(c *Config) {
	// Add flags
	config := flag.String("config", c.ConfigFile, "Path to config.")
	cport := flag.Int("cport", c.CPort, "Listen port for UPnP notify server.")
	port := flag.Int("port", c.Port, "Listen port for web server.")

	flag.Parse()

	// Add flags to config
	c.CPort = *cport
	c.ConfigFile = *config
	c.Port = *port
}
