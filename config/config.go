package config

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/ItsNotGoodName/go-upnpsub"
)

type Config struct {
	APIURI         string
	CPort          int
	CPortFlag      bool
	ConfigFile     string
	ConfigFileFlag bool
	Presets        []Preset
	Port           int
	PortFlag       bool
	PresetsEnabled bool
}

type ConfigJSON struct {
	Port    int      `json:"port"`
	CPort   int      `json:"cport"`
	Presets []Preset `json:"presets"`
}

type Preset struct {
	URL     string `json:"url"`
	NewName string `json:"newName"`
	NewURL  string `json:"newUrl"`
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

	// Validate urls
	var urls []string
	for i := range c.Presets {
		urls = append(urls, c.Presets[i].URL)
	}
	if err := ValidateURLS(urls); err != nil {
		log.Fatal("config.NewConfig:", err)
	}

	c.PresetsEnabled = len(c.Presets) > 0

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

	// Set flags in config
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "config" {
			c.ConfigFileFlag = true
		}
		if f.Name == "cport" {
			c.CPortFlag = true
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
	c.Presets = cj.Presets

	// Use config file if not using flags
	if !c.CPortFlag && cj.CPort != 0 {
		c.CPort = cj.CPort
	}
	if !c.PortFlag && cj.Port != 0 {
		c.Port = cj.Port
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
