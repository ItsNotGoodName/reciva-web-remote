package radio

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

// Hub handles creating Radios and pushing State changes to HubClients.
type Hub struct {
	cp            *upnpsub.ControlPoint           // cp is used to create subscriptions.
	discoverChan  chan chan error                 // discoverChan is used to discover radios.
	stateOPS      chan func(map[*chan State]bool) // stateOPS is used to push State changes to HubClients.
	stopChan      chan chan error                 // buryChan is used to remove all radios.
	PresetMutator func(context.Context, *Preset)  // PresetMutator is used to change presets.

	radiosMu sync.RWMutex      // radiosMu is used to protect radios map.
	radios   map[string]*Radio // radios is used to store all the Radios.
}

// Radio represents the radio on the network.
type Radio struct {
	Cancel           context.CancelFunc    // Cancel should be called when the radio is no longer needed.
	Client           goupnp.ServiceClient  // Client is the SOAP client.
	Subscription     *upnpsub.Subscription // Subscription that belongs to this Radio.
	UUID             string                // UUID is unique and will not change after it has been set.
	refreshPresets   chan bool             // RefreshPresets is used to refresh presets.
	ctx              context.Context       // ctx is the context for radioLoop.
	getStateChan     chan State            // getStateChan returns a copy of the current State.
	h                *Hub                  // h is the hub that this radio belongs to.
	state            *State                // state represents the current State of the Radio.
	updateVolumeChan chan int              // updateVolumeChan is used to update State's volume.
}

type Preset struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
	Title  string `json:"-"`
	URL    string `json:"url"`
}

// State of the radio.
type State struct {
	IsMuted    *bool    `json:"isMuted,omitempty"`    // IsMuted represents if the Radio's volume is muted.
	Metadata   *string  `json:"metadata,omitempty"`   // Metadata that is received from the stream url.
	Name       string   `json:"name,omitempty"`       // Name of the radio.
	NumPresets int      `json:"numPresets,omitempty"` // NumPresets on the radio, it will not change after it is set.
	Power      *bool    `json:"power,omitempty"`      // Power represents if the radio is not in standby.
	Presets    []Preset `json:"presets,omitempty"`    // Presets on the radio, it's length should not change after it is set.
	Preset     int      `json:"preset,omitempty"`     // Preset is the current preset that is playing, -1 means it is unknown.
	State      string   `json:"state,omitempty"`      // State is either playing, connecting, or stopped.
	Title      string   `json:"title,omitempty"`      // Title of the preset.
	UUID       string   `json:"uuid"`                 // UUID will not change once it has been set.
	URL        string   `json:"url,omitempty"`        // URL of the stream that is currently selected.
	Volume     *int     `json:"volume,omitempty"`     // Volume of the Radio.
}

type stateXML struct {
	Metadata string `xml:"playback-details>stream>metadata"`
	State    string `xml:"playback-details>state"`
	Title    string `xml:"playback-details>stream>title"`
	URL      string `xml:"playback-details>stream>url"`
}
