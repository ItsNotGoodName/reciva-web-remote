package config

import (
	"flag"

	"github.com/ItsNotGoodName/goupnpsub"
)

type Config struct {
	APIURI     string
	CPort      int
	CPortFlag  bool
	ConfigPath string
	Port       int
	PortFlag   bool
}

const (
	APIURI      = "/v1"
	DefaultPort = 8080
)

func NewConfig(options ...func(*Config)) *Config {
	c := &Config{
		APIURI: APIURI,
		CPort:  goupnpsub.DefaultPort,
		Port:   DefaultPort,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func WithFlag(c *Config) {
	cport := flag.Int("cport", c.CPort, "Listen port for UPnP notify server.")
	port := flag.Int("port", c.Port, "Listen port for web server.")
	config := flag.String("config", c.ConfigPath, "Path to config location.")

	flag.Parse()

	c.CPort = *cport
	c.Port = *port
	c.ConfigPath = *config

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			c.CPortFlag = true
		}
		if f.Name == "cport" {
			c.CPortFlag = true
		}
	})
}
