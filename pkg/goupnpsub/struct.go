package goupnpsub

import (
	"encoding/xml"
	"sync"
)

// ControlPoint handles the HTTP notify server and keeps track of subscriptions.
type ControlPoint struct {
	listenUri     string                   // The URI that the publisher sends it's notify request.
	listenPort    string                   // The port that the publisher sends it's notify request.
	sidMap        map[string]*Subscription // SID to Subscription map.
	sidMapRWMutex sync.RWMutex             // Read and write lock for sidMap.
}

// Subscription represents the subscription to the UPnP publisher.
type Subscription struct {
	EventChan   chan *Event // Events sent from publisher to this Subscription.
	ActiveChan  chan bool   // ActiveChan returns true if the Subscription is active.
	renewChan   chan bool   // Force a Subscription renewal.
	eventUrl    string      // URL to send subscribe request.
	callbackUrl string      // The full url that the publisher sends it's notify request.
	sid         string      // Do not read or write unless in subsciptionLoop.
	timeout     int         // Do not read or write unless in subscriptionLoop.
}

// Property are notify request's property.
type Property struct {
	Name  string // Name of inner field from UPnP property.
	Value string // Value of inner field from UPnP property.
}

// Event represents a parsed notify request.
type Event struct {
	// TODO: Add seq number
	sid        string
	Properties []Property
}

// propertyVariableXML represents the inner information of the property tag in the notify request's xml.
type propertyVariableXML struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// propertyXML represents property tag in the notify request's xml.
type propertyXML struct {
	Property propertyVariableXML `xml:",any"`
}

// eventXML represents a notify request's xml.
type eventXML struct {
	XMLName    xml.Name      `xml:"urn:schemas-upnp-org:event-1-0 propertyset"`
	Properties []propertyXML `xml:"urn:schemas-upnp-org:event-1-0 property"`
}
