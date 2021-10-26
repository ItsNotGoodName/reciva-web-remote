package config

import (
	"flag"
	"log"
	"strings"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

type Config struct {
	ConfigPath    string
	Port          int
	PortFlag      bool
	CPort         int
	CPortFlag     bool
	EnablePresets bool
	Presets       []string
	APIURI        string
}

const DefaultPort = 8080
const DefaultFile = "reciva-web-remote.json"

func NewConfig() *Config {
	port := flag.Int("port", DefaultPort, "Listen port for web server.")
	portFlag := false
	cport := flag.Int("cport", goupnpsub.DefaultPort, "Listen port for UPnP notify server.")
	cportFlag := false
	presetsFlag := flag.String("presets", "", "List of presets to host seperated by comma (ex. /01.m3u,/02.m3u).")
	config := flag.String("config", DefaultFile, "Path to config location.")
	enablePresets := false
	var presets []string

	flag.Parse()

	// Check if flag was specified for port and cport
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			portFlag = true
		}
		if f.Name == "cport" {
			cportFlag = true
		}
		log.Println(f.Name)
	})

	// Enable presets based on presetsFlag
	if *presetsFlag == "" {
		presets = make([]string, 0)
	} else {
		presets = strings.Split(*presetsFlag, ",")
		enablePresets = true
	}

	return &Config{
		Port:          *port,
		PortFlag:      portFlag,
		CPort:         *cport,
		CPortFlag:     cportFlag,
		EnablePresets: enablePresets,
		Presets:       presets,
		ConfigPath:    *config,
		APIURI:        "/v1",
	}
}
