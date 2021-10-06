package radio

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/huin/goupnp"
)

// Hub handles creating Radios and pushing State changes to HubClients.
type Hub struct {
	Register   chan *HubClient         // Register requests from clients.
	Unregister chan *HubClient         // Unregister requests from clients.
	clients    map[*HubClient]bool     // clients are State receivers.
	cp         *goupnpsub.ControlPoint // cp is used to create subscriptions.
	stateChan  chan State              // stateChan gets State changes from radioLoop.
}

// HubClient receives State from Hub.
type HubClient struct {
	Send chan State // Send receives State from Hub.
}

// Radio represents the radio on the network.
type Radio struct {
	Cancel           context.CancelFunc      // Cancel should be called when the radio is being removed.
	Client           *Client                 // Client that belongs to this Radio.
	GetStateChan     chan State              // GetStateChan returns a copy of the current State.
	UpdateVolumeChan chan int                // UpdateVolumeChan is used to update State's volume.
	state            *State                  // state represents the current State of the Radio.
	stateChan        chan<- State            // stateChan is written to when State changes.
	subscription     *goupnpsub.Subscription // Subscription that belongs to this Radio.
}

// State of the radio.
type State struct {
	IsMuted  bool   `json:"isMuted"`                                         // IsMuted represents if the Radio's volume is muted.
	Metadata string `json:"metadata" xml:"playback-details>stream>metadata"` // Metadata that is received from the stream url.
	Power    bool   `json:"power"`                                           // Power represents if the radio is not in standby.
	Presets  int    `json:"presets"`                                         // Presets represents the max amount of presets and starts at 1.
	State    string `json:"state" xml:"playback-details>state"`              // State is either playing, connecting, or stopped.
	Title    string `json:"title" xml:"playback-details>stream>title"`       // Title of the preset.
	UUID     string `json:"uuid"`                                            // UUID will not change once it has been set.
	Url      string `json:"url" xml:"playback-details>stream>url"`           // Url of the stream that is currently selected.
	Volume   int    `json:"volume"`                                          // Volume of the Radio.
}

// Client is used to send UPnP commands to the radio.
type Client struct {
	goupnp.ServiceClient
	UUID string
}
