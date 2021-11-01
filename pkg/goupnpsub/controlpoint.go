package goupnpsub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// NewControlPoint creates a ControlPoint.
func NewControlPoint() *ControlPoint {
	return NewControlPointWithPort(DefaultPort)
}

// NewControlPoint creates a ControlPoint that listens on a specific port.
func NewControlPointWithPort(listenPort int) *ControlPoint {
	cp := &ControlPoint{
		listenURI:  ListenURI,
		listenPort: fmt.Sprint(listenPort),
		sidMap:     make(map[string]*Subscription),
		sidMapRWMu: sync.RWMutex{},
	}

	http.Handle(cp.listenURI, cp)
	return cp
}

// Start HTTP server that listens for notify requests.
func (cp *ControlPoint) Start() {
	log.Println("ControlPoint.Start: listening on port", cp.listenPort)
	log.Fatal(http.ListenAndServe(":"+cp.listenPort, nil))
}

// NewSubscription creates and returns a Subscription.
func (cp *ControlPoint) NewSubscription(ctx context.Context, eventURL *url.URL) (*Subscription, error) {
	// Find callback ip
	callbackIP, err := findCallbackIP(eventURL)
	if err != nil {
		return nil, err
	}

	// Create sub
	sub := &Subscription{
		Active:        make(chan bool),
		Done:          make(chan bool),
		Event:         make(chan *Event),
		callback:      "<http://" + callbackIP + ":" + cp.listenPort + cp.listenURI + ">",
		eventURL:      eventURL.String(),
		renewChan:     make(chan bool),
		setActiveChan: make(chan bool),
	}

	// Start sub loops
	go cp.subscriptionLoop(ctx, sub)
	go sub.activeLoop()

	return sub, nil
}

// ServeHTTP handles notify requests from event publishers.
func (cp *ControlPoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get NT and NTS
	nt, nts := r.Header.Get("NT"), r.Header.Get("NTS")
	if nt == "" || nts == "" {
		log.Println("ControlPoint.ServeHTTP(WARNING): request has no nt or nts")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get SEQ
	var seq int
	if seqStr := r.Header.Get("SEQ"); seqStr != "" {
		seqInt, err := strconv.Atoi(seqStr)
		if err != nil {
			log.Println("ControlPoint.ServeHTTP(WARNING): invalid seq", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		seq = seqInt
	}

	// Validate NT and NTS
	if nt != NT || nts != NTS {
		log.Printf("ControlPoint.ServeHTTP(WARNING): invalid nt or nts, %s, %s", nt, nts)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// Get SID
	sid := r.Header.Get("SID")

	// Find sub from sidMap using SID
	cp.sidMapRWMu.RLock()
	sub, ok := cp.sidMap[sid]
	cp.sidMapRWMu.RUnlock()
	if !ok {
		log.Println("ControlPoint.ServeHTTP(WARNING): sid not found or valid,", sid)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// Get body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ControlPoint.ServeHTTP:", err)
		return
	}

	// Parse xmlEvent from body
	xmlEvent, err := unmarshalEventXML(body)
	if err != nil {
		log.Println("ControlPoint.ServeHTTP:", err)
		return
	}

	// Parse properties from xmlEvent
	properties := unmarshalProperties(xmlEvent)

	// Try to send event to sub's Event
	t := time.NewTimer(DefaultTimeout)
	select {
	case <-t.C:
		log.Println("ControlPoint.ServeHTTP(ERROR): could not send event to subscription's Event")
	case sub.Event <- &Event{Properties: properties, SEQ: seq, sid: sid}:
		if !t.Stop() {
			<-t.C
		}
	}
}

// subscribe sends SUBSCRIBE request to event publisher.
func (cp *ControlPoint) subscribe(ctx context.Context, sub *Subscription) error {
	// Create request
	req, err := http.NewRequest("SUBSCRIBE", sub.eventURL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// Add headers to request
	req.Header.Add("CALLBACK", sub.callback)
	req.Header.Add("NT", NT)
	req.Header.Add("TIMEOUT", Timeout)

	// Execute request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check if request failed
	if res.StatusCode != http.StatusOK {
		return errors.New("invalid status " + res.Status)
	}

	// Get SID
	sid := res.Header.Get("sid")
	if sid == "" {
		return errors.New("subscribe's response has no sid")
	}

	cp.sidMapRWMu.Lock()

	// Delete old SID to sub mapping
	delete(cp.sidMap, sub.sid)

	// Add new SID to sub mapping
	sub.sid = sid
	cp.sidMap[sid] = sub

	cp.sidMapRWMu.Unlock()

	// Update sub's timeout
	timeout, err := unmarshalTimeout(res.Header.Get("timeout"))
	if err != nil {
		return err
	}
	sub.timeout = timeout

	return nil
}

// subscriptionLoop handles sending subscribe requests to event publisher.
func (cp *ControlPoint) subscriptionLoop(ctx context.Context, sub *Subscription) {
	log.Println("ControlPoint.subscriptionLoop: started")

	defer close(sub.Done)

	// Subscribe
	t := time.NewTimer(cp.renew(ctx, sub))

	for {
		select {
		case <-ctx.Done():
			log.Println("ControlPoint.subscriptionLoop: ctx is done, starting cleanup")

			// Delete sub.sid from sidMap
			cp.sidMapRWMu.Lock()
			delete(cp.sidMap, sub.sid)
			cp.sidMapRWMu.Unlock()

			// Unsubscribe
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
			if err := sub.unsubscribe(ctx); err != nil {
				log.Print("ControlPoint.subscriptionLoop:", err)
			}
			cancel()

			log.Println("ControlPoint.subscriptionLoop: cleanup finished")
			return
		case <-sub.renewChan:
			log.Println("ControlPoint.subscriptionLoop: starting manual renewal")

			// Manual renew
			if !t.Stop() {
				<-t.C
			}
			t.Reset(cp.renew(ctx, sub))
		case <-t.C:
			// Renew
			t.Reset(cp.renew(ctx, sub))
		}
	}
}

// renew handles subscribing or resubscribing.
func (cp *ControlPoint) renew(ctx context.Context, sub *Subscription) time.Duration {
	if !<-sub.Active {
		if err := cp.subscribe(ctx, sub); err != nil {
			log.Print("ControlPoint.subscriptionLoop:", err)
			return getRenewDuration(sub)
		}
		sub.setActive(ctx, true)
		duration := getRenewDuration(sub)
		log.Printf("ControlPoint.subscriptionLoop: subscribe successful, will resubscribe in %s intervals", duration)
		return duration
	}
	if err := sub.resubscribe(ctx); err != nil {
		sub.setActive(ctx, false)
		log.Print("ControlPoint.subscriptionLoop:", err)
		return DefaultTimeout
	}
	return getRenewDuration(sub)
}
