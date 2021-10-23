package config

import (
	"flag"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

type Config struct {
	Port  int
	CPort int
}

const DefaultPort = 8080

func NewConfig() *Config {
	port := flag.Int("port", DefaultPort, "Listen port for web server.")
	cport := flag.Int("cport", goupnpsub.DefaultPort, "Listen port for UPnP notify server.")

	flag.Parse()

	return &Config{
		Port:  *port,
		CPort: *cport,
	}
}
