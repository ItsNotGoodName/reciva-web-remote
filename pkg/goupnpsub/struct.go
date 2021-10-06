package goupnpsub

import (
	"encoding/xml"
	"sync"
)

// ControlPoint handles the http notify server and keeps track of subscriptions.
type ControlPoint struct {
	listenUri     string                   // The URI that the publisher sends it's notify request.
	listenPort    string                   // The port that the publisher sends it's notify request.
	sidMap        map[string]*Subscription // SID to Subscription map.
	sidMapRWMutex sync.RWMutex             // Read and write lock for sidMap.
}

// Subscription represents the subscription to the upnp publisher.
type Subscription struct {
	EventChan   chan *Event // Events for this subscriptions from publisher.
	ActiveChan  chan bool   // ActiveChan returns true if the the subscription is active.
	renewChan   chan bool   // Force a subscription renewel.
	eventUrl    string      // URL to send subscribe request.
	callbackUrl string      // The full url that the publisher sends it's notify request.
	sid         string      // Do not read or write unless in subsciptionLoop.
	timeout     int         // Do not read or write unless in subscriptionLoop.
}

// Property are notify request's property.
type Property struct {
	Name  string // Name of inner field from upnp property.
	Value string // Value of inner field from upnp property.
}

// Event represents a notify request's information.
type Event struct {
	// TODO: Add seq number
	sid        string
	Properties []Property
}

// propertyVariableXML represents the inner information of property tag in notify request's xml.
type propertyVariableXML struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// propertyXML represents property tag in notify request's xml.
type propertyXML struct {
	Property propertyVariableXML `xml:",any"`
}

// eventXML represents notify request's xml.
type eventXML struct {
	XMLName    xml.Name      `xml:"urn:schemas-upnp-org:event-1-0 propertyset"`
	Properties []propertyXML `xml:"urn:schemas-upnp-org:event-1-0 property"`
}
