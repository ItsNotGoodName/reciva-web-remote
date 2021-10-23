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
	defaultTimeout = "Second-300"
	DefaultPort    = 8058
)

var timeoutReg = regexp.MustCompile(`(?i)second-([0-9]*)`)

func parseEventXML(body []byte) (*eventXML, error) {
	xmlEvent := &eventXML{}
	return xmlEvent, xml.Unmarshal(body, xmlEvent)
}

func parseProperties(xmlEvent *eventXML) []Property {
	properties := make([]Property, len(xmlEvent.Properties))
	for i := range xmlEvent.Properties {
		properties[i].Name = xmlEvent.Properties[i].Property.XMLName.Local
		properties[i].Value = xmlEvent.Properties[i].Property.Value
	}
	return properties
}

func parseTimeout(timeoutPre string) (int, error) {
	timeoutArr := timeoutReg.FindStringSubmatch(timeoutPre)
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
		return 0, errors.New("timeout is invalid" + timeoutString)
	}
	return timeoutInt, nil
}

func getRenewDuration(timeout int) time.Duration {
	return time.Duration(timeout/2) * time.Second
}

func getCallbackIP(url *url.URL) (string, error) {
	// https://stackoverflow.com/a/37382208
	conn, err := net.Dial("udp", url.Host)
	conn.Close()
	if err != nil {
		return "", err
	}
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}
