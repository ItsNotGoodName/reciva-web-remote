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
	URLS           []string
	DB             string
	DBFlag         bool
	Port           int
	PortFlag       bool
	PresetsEnabled bool
}

type Preset struct {
	URL string
	URI string
}

type ConfigJSON struct {
	DBPath string   `json:"db"`
	Port   int      `json:"port"`
	CPort  int      `json:"cport"`
	URLS   []string `json:"presets"`
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
	c.PresetsEnabled = len(c.URLS) > 0
	return c
}

func WithFlag(c *Config) {
	// Add flags
	URLS := flag.String("presets", "", "List of preset URLs, ex. 'http://example.com//01.m3u,http://example.com/02.m3u'.")
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
	if *URLS == "" {
		c.URLS = []string{}
	} else {
		URLS := strings.Split(*URLS, ",")
		err := ValidateURLS(URLS)
		if err != nil {
			log.Fatal("Config.WithFlag:", err)
		}
		c.URLS = URLS
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
	if !c.ConfigFileFlag && len(cj.URLS) > 0 {
		err := ValidateURLS(cj.URLS)
		if err != nil {
			log.Fatal("Config.WithFile:", err)
		}
		c.URLS = cj.URLS
	}
}

// ValidateURLS validates a list of preset URLs.
func ValidateURLS(urls []string) error {
	for _, u := range urls {
		_, err := url.Parse(u)
		if err != nil {
			return err
		}
	}
	return nil
}
