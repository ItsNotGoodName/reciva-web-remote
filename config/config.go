package config

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/ItsNotGoodName/go-upnpsub"
)

type Config struct {
	APIURI         string
	CPort          int
	CPortFlag      bool
	ConfigFile     string
	ConfigFileFlag bool
	URIS           []string
	DB             string
	DBFlag         bool
	Port           int
	PortFlag       bool
	PresetsEnabled bool
}

type ConfigJSON struct {
	DBPath string   `json:"db"`
	Port   int      `json:"port"`
	CPort  int      `json:"cport"`
	URIS   []string `json:"presets"`
}

const (
	APIURI      = "/v1"
	ConfigFile  = "reciva-web-remote.json"
	DBPath      = "reciva-web-remote.db"
	DefaultPort = 8080
)

func NewConfig(options ...func(*Config)) *Config {
	c := &Config{
		APIURI:     APIURI,
		CPort:      upnpsub.DefaultPort,
		ConfigFile: ConfigFile,
		DB:         DBPath,
		Port:       DefaultPort,
	}
	for _, option := range options {
		option(c)
	}
	if len(c.URIS) > 0 {
		c.PresetsEnabled = true
		if err := ValidURIS(c.URIS); err != nil {
			log.Fatal(err)
		}
	}
	return c
}

func WithFlag(c *Config) {
	// Add flags
	URIS := flag.String("presets", "", "List of preset URI, ex. '/01.m3u,/02.m3u,/02.m3u'.")
	config := flag.String("config", c.ConfigFile, "Path to config.")
	cport := flag.Int("cport", c.CPort, "Listen port for UPnP notify server.")
	db := flag.String("db", c.DB, "Path to database.")
	port := flag.Int("port", c.Port, "Listen port for web server.")

	flag.Parse()

	// Add flags to config
	c.CPort = *cport
	c.ConfigFile = *config
	c.DB = *db
	c.Port = *port
	if *URIS == "" {
		c.URIS = make([]string, 0)
	} else {
		c.URIS = strings.Split(*URIS, ",")
	}

	// Set flags in config
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "config" {
			c.ConfigFileFlag = true
		}
		if f.Name == "cport" {
			c.CPortFlag = true
		}
		if f.Name == "db" {
			c.DBFlag = true
		}
		if f.Name == "port" {
			c.PortFlag = true
		}
	})
}

func WithFile(c *Config) {
	// Check if config file exists
	if _, err := os.Stat(c.ConfigFile); errors.Is(err, os.ErrNotExist) {
		if c.ConfigFileFlag {
			log.Fatalf("Config file %s not found.", c.ConfigFile)
		}
		return
	}

	// Read config file
	data, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal config file
	cj := ConfigJSON{}
	if err := json.Unmarshal(data, &cj); err != nil {
		log.Fatal(err)
	}

	// Use config file if not using flags
	if !c.CPortFlag && cj.CPort != 0 {
		c.CPort = cj.CPort
	}
	if !c.DBFlag && cj.DBPath != "" {
		c.DB = cj.DBPath
	}
	if !c.PortFlag && cj.Port != 0 {
		c.Port = cj.Port
	}
	if !c.ConfigFileFlag && len(cj.URIS) > 0 {
		c.URIS = cj.URIS
	}
}

// ValidURIS checks if URIs are valid.
func ValidURIS(uris []string) error {
	for _, uri := range uris {
		_, err := url.ParseRequestURI(uri)
		if err != nil {
			return err
		}
	}
	return nil
}
