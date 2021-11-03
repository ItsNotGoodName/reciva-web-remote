package config

import (
	"flag"
	"strings"

	"github.com/ItsNotGoodName/go-upnpsub"
)

type Config struct {
	APIURI         string
	CPort          int
	CPortFlag      bool
	ConfigPath     string
	URIS           []string
	DBPath         string
	Port           int
	PortFlag       bool
	PresetsEnabled bool
}

const (
	APIURI      = "/v1"
	DBPath      = "reciva-web-remote.db"
	DefaultPort = 8080
)

func NewConfig(options ...func(*Config)) *Config {
	c := &Config{
		APIURI: APIURI,
		CPort:  upnpsub.DefaultPort,
		DBPath: DBPath,
		Port:   DefaultPort,
	}
	for _, option := range options {
		option(c)
	}
	if len(c.URIS) > 0 {
		c.PresetsEnabled = true
	}
	return c
}

func WithFlag(c *Config) {
	cport := flag.Int("cport", c.CPort, "Listen port for UPnP notify server.")
	port := flag.Int("port", c.Port, "Listen port for web server.")
	config := flag.String("config", c.ConfigPath, "Path to config location.")
	URIS := flag.String("presets", "", "List of preset URIs, ex. '/01.m3u'")

	flag.Parse()

	c.CPort = *cport
	c.Port = *port
	c.ConfigPath = *config
	if *URIS == "" {
		c.URIS = make([]string, 0)
	} else {
		c.URIS = strings.Split(*URIS, ",")
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			c.CPortFlag = true
		}
		if f.Name == "cport" {
			c.CPortFlag = true
		}
	})
}
