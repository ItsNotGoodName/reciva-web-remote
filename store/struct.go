package store

import (
	"context"
	"sync"
)

type Settings struct {
	CPort   int      `json:"cport"`
	Port    int      `json:"port"`
	Presets []Preset `json:"presets"`
	Streams []Stream `json:"streams"`
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
	dctx           context.Context
	file           string
	queueWriteChan chan bool
	readChan       chan chan error
	sg             *Settings
	sgMutex        sync.Mutex
	writeChan      chan chan error
}
