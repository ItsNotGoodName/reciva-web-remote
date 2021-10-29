package radio

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/huin/goupnp"
)

// Hub handles creating Radios and pushing State changes to HubClients.
type Hub struct {
	Register         chan *chan State        // Register requests from clients.
	Unregister       chan *chan State        // Unregister requests from clients.
	clients          map[*chan State]bool    // clients are registered to receive state changes from all radios.
	cp               *goupnpsub.ControlPoint // cp is used to create subscriptions.
	receiveStateChan chan State              // receiveStateChan gets State changes from radioLoop.
}

// Radio represents the radio on the network.
type Radio struct {
	Cancel            context.CancelFunc      // Cancel should be called when the radio is no longer needed.
	Client            goupnp.ServiceClient    // Client is the SOAP client.
	Subscription      *goupnpsub.Subscription // Subscription that belongs to this Radio.
	UUID              string                  // UUID is unique and will not change after it has been set.
	dctx              context.Context         // dctx is the context for radioLoop.
	getStateChan      chan State              // getStateChan returns a copy of the current State.
	sendStateChan     chan<- State            // sendStateChan sends state to hubLoop.
	state             *State                  // state represents the current State of the Radio.
	updatePresetsChan chan []Preset           // updatePresetsChan is used to update State's presets.
	updateVolumeChan  chan int                // updateVolumeChan is used to update State's volume.
}

type Preset struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// State of the radio.
type State struct {
	IsMuted    *bool    `json:"isMuted,omitempty"`    // IsMuted represents if the Radio's volume is muted.
	Metadata   *string  `json:"metadata,omitempty"`   // Metadata that is received from the stream url.
	Name       string   `json:"name,omitempty"`       // Name of the radio.
	NumPresets int      `json:"numPresets,omitempty"` // NumPresets on the radio, it will not change after it is set.
	Power      *bool    `json:"power,omitempty"`      // Power represents if the radio is not in standby.
	Presets    []Preset `json:"presets,omitempty"`    // Presets on the radio.
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
