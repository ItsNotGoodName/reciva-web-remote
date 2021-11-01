package goupnpsub

import (
	"encoding/xml"
	"sync"
)

// ControlPoint handles the HTTP notify server and keeps track of subscriptions.
type ControlPoint struct {
	listenURI  string // listenURI is the URI for consuming notify requests.
	listenPort string // listenPort is the HTTP server's port.

	sidMapRWMu sync.RWMutex             // sidMapRWMu read-write mutex.
	sidMap     map[string]*Subscription // sidMap hold all active subscriptions.
}

// Subscription represents the subscription to the UPnP publisher.
type Subscription struct {
	Active        chan bool   // Active represents the subscription status to publisher.
	Done          chan bool   // Done is closed when subscriptionLoop is finished.
	Event         chan *Event // Event is the UPnP events from publisher.
	callback      string      // callback is part of the UPnP header.
	eventURL      string      // eventURL is the URL we are subscribed on the event publisher.
	renewChan     chan bool   // renewChan forces a subscription renewal.
	setActiveChan chan bool   // SetActiveChan sets the active status.
	sid           string      // sid the header set by the UPnP publisher.
	timeout       int         // timeout is the timeout seconds received from UPnP publisher.
}

// Property is the notify request's property.
type Property struct {
	Name  string // Name of inner field from UPnP property.
	Value string // Value of inner field from UPnP property.
}

// Event represents a parsed notify request.
type Event struct {
	Properties []Property
	SEQ        int
	sid        string
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
