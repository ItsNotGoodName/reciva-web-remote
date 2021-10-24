package store

import (
	"context"
	"sync"
)

type Settings struct {
	Port    int      `json:"port"`
	CPort   int      `json:"cport"`
	Streams []Stream `json:"streams"`
	Presets []Preset `json:"presets"`
}

type Stream struct {
	SID     int    `json:"sid"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Preset struct {
	URI string `json:"uri"`
	SID int    `json:"sid"`
}

type Store struct {
	Cancel         context.CancelFunc
	file           string
	dctx           context.Context
	st             *Settings
	stMutex        sync.Mutex
	writeChan      chan chan error
	readChan       chan chan error
	queueWriteChan chan bool
}
