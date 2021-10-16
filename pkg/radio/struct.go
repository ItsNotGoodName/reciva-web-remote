package radio

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/huin/goupnp"
)

// Hub handles creating Radios and pushing State changes to HubClients.
type Hub struct {
	Register     chan *HubClient         // Register requests from clients.
	Unregister   chan *HubClient         // Unregister requests from clients.
	clients      map[*HubClient]bool     // clients are registered to receive state changes from all radios.
	cp           *goupnpsub.ControlPoint // cp is used to create subscriptions.
	allStateChan chan State              // allStateChan gets State changes from radioLoop.
}

// HubClient receives State from Hub.
type HubClient struct {
	Send chan State // Send receives State from Hub.
}

// Radio represents the radio on the network.
type Radio struct {
	Cancel           context.CancelFunc      // Cancel should be called when the radio is being removed.
	Client           goupnp.ServiceClient    // Client is the goupnp SOAP client.
	Subscription     *goupnpsub.Subscription // Subscription that belongs to this Radio.
	UUID             string                  // UUID is unique and should not change after it has been set.
	allStateChan     chan<- State            // allStateChan is written to when State changes.
	dctx             context.Context         // dctx is the context for radioLoop.
	getStateChan     chan State              // GetStateChan returns a copy of the current State.
	state            *State                  // state represents the current State of the Radio.
	updateVolumeChan chan int                // UpdateVolumeChan is used to update State's volume.
}

// State of the radio.
type State struct {
	IsMuted  bool   `json:"isMuted"`                                         // IsMuted represents if the Radio's volume is muted.
	Metadata string `json:"metadata" xml:"playback-details>stream>metadata"` // Metadata that is received from the stream url.
	Name     string `json:"name"`                                            // Name of the radio.
	Power    bool   `json:"power"`                                           // Power represents if the radio is not in standby.
	Presets  int    `json:"presets"`                                         // Presets represents the max amount of presets and starts at 1.
	State    string `json:"state" xml:"playback-details>state"`              // State is either playing, connecting, or stopped.
	Title    string `json:"title" xml:"playback-details>stream>title"`       // Title of the preset.
	UUID     string `json:"uuid"`                                            // UUID will not change once it has been set.
	Url      string `json:"url" xml:"playback-details>stream>url"`           // Url of the stream that is currently selected.
	Volume   int    `json:"volume"`                                          // Volume of the Radio.
}
