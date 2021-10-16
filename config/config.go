package config

import "flag"

type Config struct {
	Port  int
	CPort int
}

func NewConfig() *Config {
	port := flag.Int("port", 8080, "Listen port for web server.")
	cport := flag.Int("cport", 8058, "Listen port for UPnP notify server.")

	flag.Parse()

	return &Config{
		Port:  *port,
		CPort: *cport,
	}
}
