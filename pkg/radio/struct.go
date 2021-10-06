package radio

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/huin/goupnp"
)

// Hub handles pushing State changes to HubClients and creating Radios.
type Hub struct {
	Register   chan *HubClient         // Register requests from clients.
	Unregister chan *HubClient         // Unregister requests from clients.
	clients    map[*HubClient]bool     // clients are State receivers.
	cp         *goupnpsub.ControlPoint // cp is used to create subscriptions.
	stateChan  chan State              // stateChan gets State changes from radioLoop.
}

// HubClient recieve State from Hub.
type HubClient struct {
	Send chan State // Send receives State from Hub.
}

// Radio represents a radio.
type Radio struct {
	Cancel           context.CancelFunc      // Cancel should be called when the Radio is being removed.
	Client           *Client                 // Client is the Client that belongs to this Radio.
	GetStateChan     chan State              // GetStateChan is a channel that returns a copy of the current state.
	UpdateVolumeChan chan int                // UpdateVolumeChan is a channel that is used to update volume state.
	state            *State                  // state represents the current state of the Radio.
	stateChan        chan<- State            // uuidChan is written to when state changes.
	subscription     *goupnpsub.Subscription // Subscription is the upnp notify subscription that belongs to this Radio.
}

// State represents the current state of the radio.
type State struct {
	IsMuted  bool   `json:"isMuted"`                                         // IsMuted represents if the radio's volume is muted.
	Metadata string `json:"metadata" xml:"playback-details>stream>metadata"` // Metadata is the metadata received from the stream url.
	Power    bool   `json:"power"`                                           // Power represents if the radio is not in standby.
	Presets  int    `json:"presets"`                                         // Presets represents the max amount of presets and starts at 1.
	State    string `json:"state" xml:"playback-details>state"`              // State is either playing, connecting, or stopped.
	Title    string `json:"title" xml:"playback-details>stream>title"`       // Title is the title of the preset.
	UUID     string `json:"uuid"`                                            // UUID will not change once it has been set.
	Url      string `json:"url" xml:"playback-details>stream>url"`           // Url is the url of the stream that is currently playing.
	Volume   int    `json:"volume"`                                          // Presets represents the max amount of presets and starts at 1.
}

// Client is used to send upnp comments to the radio.
type Client struct {
	goupnp.ServiceClient
	UUID string
}
