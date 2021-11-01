package goupnpsub

import (
	"encoding/xml"
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultPort    = 8058              // DefaultPort is the port that the HTTP server listens .
	DefaultTimeout = 20 * time.Second  // DefaultTimeout is how long to wait for executing actions.
	ListenURI      = "/eventSub"       // ListenURI is the default URI that the UPnP sends notify requests.
	NT             = "upnp:event"      // NT is part of the UPnP header.
	NTS            = "upnp:propchange" // NTS is part of the UPnP header.
	Timeout        = "Second-300"      // Timeout is part of the UPnP header.
)

var timeoutReg = regexp.MustCompile(`(?i)second-([0-9]*)`)

func unmarshalEventXML(body []byte) (*eventXML, error) {
	xmlEvent := &eventXML{}
	return xmlEvent, xml.Unmarshal(body, xmlEvent)
}

// getRenewDuration returns half the sub timeout.
func getRenewDuration(sub *Subscription) time.Duration {
	return time.Duration(sub.timeout/2) * time.Second
}

func unmarshalProperties(xmlEvent *eventXML) []Property {
	properties := make([]Property, len(xmlEvent.Properties))
	for i := range xmlEvent.Properties {
		properties[i].Name = xmlEvent.Properties[i].Property.XMLName.Local
		properties[i].Value = xmlEvent.Properties[i].Property.Value
	}
	return properties
}

func unmarshalTimeout(timeout string) (int, error) {
	timeoutArr := timeoutReg.FindStringSubmatch(timeout)
	if len(timeoutArr) != 2 {
		return 0, errors.New("timeout not found")
	}
	timeoutString := timeoutArr[1]
	if strings.ToLower(timeoutString) == "infinite" {
		return 300, nil
	}
	timeoutInt, err := strconv.Atoi(timeoutString)
	if err != nil {
		return 0, err
	}
	if timeoutInt < 0 {
		return 0, errors.New("timeout is invalid, " + timeoutString)
	}
	return timeoutInt, nil
}

func findCallbackIP(url *url.URL) (string, error) {
	// https://stackoverflow.com/a/37382208
	conn, err := net.Dial("udp", url.Host)
	conn.Close()
	if err != nil {
		return "", err
	}
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}
