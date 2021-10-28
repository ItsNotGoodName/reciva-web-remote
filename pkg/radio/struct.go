package radio

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/huin/goupnp"
)

// Hub handles creating Radios and pushing State changes to HubClients.
type Hub struct {
	Register        chan *chan State        // Register requests from clients.
	Unregister      chan *chan State        // Unregister requests from clients.
	clients         map[*chan State]bool    // clients are registered to receive state changes from all radios.
	cp              *goupnpsub.ControlPoint // cp is used to create subscriptions.
	stateUpdateChan chan State              // allStateChan gets State changes from radioLoop.
}

// Radio represents the radio on the network.
type Radio struct {
	Cancel            context.CancelFunc      // Cancel should be called when the radio is being removed.
	Client            goupnp.ServiceClient    // Client is the goupnp SOAP client.
	Subscription      *goupnpsub.Subscription // Subscription that belongs to this Radio.
	UUID              string                  // UUID is unique and should not change after it has been set.
	allStateChan      chan<- State            // allStateChan is written to when State changes.
	dctx              context.Context         // dctx is the context for radioLoop.
	getStateChan      chan State              // GetStateChan returns a copy of the current State.
	state             *State                  // state represents the current State of the Radio.
	updatePresetsChan chan []Preset           // UpdateVolumeChan is used to update State's presets.
	updateVolumeChan  chan int                // UpdateVolumeChan is used to update State's volume.
}

type Preset struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	URL    string `json:"url"`
}

// State of the radio.
type State struct {
	IsMuted    bool     `json:"isMuted"`                                         // IsMuted represents if the Radio's volume is muted.
	Metadata   string   `json:"metadata" xml:"playback-details>stream>metadata"` // Metadata that is received from the stream url.
	Name       string   `json:"name"`                                            // Name of the radio.
	NumPresets int      `json:"numPresets"`                                      // NumPresets on the radio, does not change after it is set.
	Power      bool     `json:"power"`                                           // Power represents if the radio is not in standby.
	Presets    []Preset `json:"presets"`                                         // Presets on the radio.
	State      string   `json:"state" xml:"playback-details>state"`              // State is either playing, connecting, or stopped.
	Title      string   `json:"title" xml:"playback-details>stream>title"`       // Title of the preset.
	UUID       string   `json:"uuid"`                                            // UUID will not change once it has been set.
	Url        string   `json:"url" xml:"playback-details>stream>url"`           // Url of the stream that is currently selected.
	Volume     int      `json:"volume"`                                          // Volume of the Radio.
}
