package store

import "context"

type Settings struct {
	Port    int      `json:"port"`
	CPort   int      `json:"cport"`
	Streams []Stream `json:"streams"`
	Presets []Preset `json:"presets"`
}

type Stream struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Preset struct {
	URI      string `json:"uri"`
	StreamID int    `json:"id"`
}

type Store struct {
	Cancel            context.CancelFunc
	file              string
	ctx               context.Context
	getSettingsChan   chan Settings
	writeSettingsChan chan chan error
	readSettingsChan  chan chan error
	deleteStreamChan  chan int
	updateStreamChan  chan Stream
	updatePresetChan  chan Preset
	setPresetsChan    chan []Preset
}
