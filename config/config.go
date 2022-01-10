package config

import (
	"flag"
	"strconv"

	"github.com/ItsNotGoodName/go-upnpsub"
)

type Config struct {
	CPort       int
	ConfigFile  string
	Port        int
	PortStr     string
	ShowInfo    bool
	ShowVersion bool
}

const (
	ConfigFile  = "reciva-web-remote.json"
	DefaultPort = 8080
)

func NewConfig(options ...func(*Config)) *Config {
	c := &Config{
		CPort:      upnpsub.DefaultPort,
		ConfigFile: ConfigFile,
		Port:       DefaultPort,
	}

	for _, option := range options {
		option(c)
	}

	c.PortStr = strconv.Itoa(c.Port)

	return c
}

func WithFlag(c *Config) {
	// Add flags
	config := flag.String("config", c.ConfigFile, "Path to config.")
	cport := flag.Int("cport", c.CPort, "Listen port for UPnP notify server.")
	port := flag.Int("port", c.Port, "Listen port for web server.")
	showInfo := flag.Bool("info", c.ShowInfo, "Show build information about this binary.")
	ShowVersion := flag.Bool("version", c.ShowVersion, "Show version.")

	flag.Parse()

	// Add flags to config
	c.CPort = *cport
	c.ConfigFile = *config
	c.Port = *port
	c.ShowInfo = *showInfo
	c.ShowVersion = *ShowVersion
}
