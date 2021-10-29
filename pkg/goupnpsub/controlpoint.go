package goupnpsub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
		listenURI:     ListenURI,
		listenPort:    fmt.Sprint(listenPort),
		sidMap:        make(map[string]*Subscription),
		sidMapRWMutex: sync.RWMutex{},
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
		EventChan:     make(chan *Event, 10),
		GetActiveChan: make(chan bool),
		callbackURL:   "<http://" + callbackIP + ":" + cp.listenPort + cp.listenURI + ">",
		eventURL:      eventURL.String(),
		renewChan:     make(chan bool),
		setActiveChan: make(chan bool),
	}

	// Start sub loops
	go cp.subscriptionLoop(ctx, sub)
	go sub.activeLoop(ctx)

	return sub, nil
}

// ServeHTTP handles notify requests from event publishers.
func (cp *ControlPoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate NT and NTS
	nt, nts := r.Header.Get("NT"), r.Header.Get("NTS")
	if nt == "" || nts == "" {
		log.Println("ControlPoint.ServeHTTP(WARNING): request has no nt or nts")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if nt != NT || nts != NTS {
		log.Printf("ControlPoint.ServeHTTP(WARNING): bad nt or nts, %s, %s", nt, nts)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// Validate SID
	sid := r.Header.Get("SID")
	if sid == "" {
		log.Println("ControlPoint.ServeHTTP(WARNING): request has no sid")
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// Find sub from sidMap using SID
	cp.sidMapRWMutex.RLock()
	sub, ok := cp.sidMap[sid]
	cp.sidMapRWMutex.RUnlock()
	if !ok {
		log.Println("ControlPoint.ServeHTTP(WARNING): could not find sid in sidMap")
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
	xmlEvent, err := parseEventXML(body)
	if err != nil {
		log.Println("ControlPoint.ServeHTTP:", err)
		return
	}

	// Parse properties from xmlEvent
	properties := parseProperties(xmlEvent)

	// Create event using SID and properties
	event := Event{sid: sid, Properties: properties}

	// Try to send event to sub's EventChan, fail after waiting for 20 seconds
	select {
	case <-time.After(20 * time.Second):
		log.Println("ControlPoint.ServeHTTP(ERROR): could not send event to subscription's EventChan")
	case sub.EventChan <- &event:
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
	req.Header.Add("CALLBACK", sub.callbackURL)
	req.Header.Add("NT", NT)
	req.Header.Add("TIMEOUT", DefaultTimeout)

	// Execute request
	client := http.Client{}
	res, err := client.Do(req)

	// Lock map as soon as possible to prevent race condition with ServeHTTP
	cp.sidMapRWMutex.Lock()

	// Check request and get SID from request
	if err != nil {
		cp.sidMapRWMutex.Unlock()
		return err
	}
	if res.StatusCode != http.StatusOK {
		cp.sidMapRWMutex.Unlock()
		return errors.New("invalid status " + res.Status)
	}
	sid := res.Header.Get("sid")
	if sid == "" {
		cp.sidMapRWMutex.Unlock()
		return errors.New("subscribe's response has no sid")
	}

	// Delete old SID of sub in map and update with new SID from request
	delete(cp.sidMap, sub.sid)
	sub.sid = sid
	cp.sidMap[sid] = sub

	// Unlock map
	cp.sidMapRWMutex.Unlock()

	// Update sub's timeout with request's SID
	timeout, err := parseTimeout(res.Header.Get("timeout"))
	if err != nil {
		return err
	}
	sub.timeout = timeout

	return nil
}

// subscriptionLoop handles subscribing and renewing subscriptions.
func (cp *ControlPoint) subscriptionLoop(ctx context.Context, sub *Subscription) {
	log.Println("ControlPoint.subscriptionLoop: started")

	// Renew sub and get d til next renewal
	t := time.NewTimer(cp.renew(ctx, sub))

	for {
		select {
		case <-ctx.Done():
			log.Println("ControlPoint.subscriptionLoop: ctx is done, starting cleanup")

			// Delete sub.sid from sidMap
			cp.sidMapRWMutex.Lock()
			delete(cp.sidMap, sub.sid)
			cp.sidMapRWMutex.Unlock()

			// Unsubscribe
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if err := sub.unsubscribe(ctx); err != nil {
				log.Print("ControlPoint.subscriptionLoop:", err)
			}

			log.Println("ControlPoint.subscriptionLoop: cleanup finished")
			return
		case <-sub.renewChan:
			log.Println("ControlPoint.subscriptionLoop: renewChan received")
			if !t.Stop() {
				<-t.C
			}
			t.Reset(cp.renew(ctx, sub))
		case <-t.C:
			t.Reset(cp.renew(ctx, sub))
		}
	}
}

// renew handles subscribing or resubscribing.
func (cp *ControlPoint) renew(ctx context.Context, sub *Subscription) time.Duration {
	if !<-sub.GetActiveChan {
		if err := cp.subscribe(ctx, sub); err != nil {
			log.Print("ControlPoint.subscriptionLoop:", err)
			return sub.getRenewDuration()
		}
		sub.setActiveChan <- true
		duration := sub.getRenewDuration()
		log.Printf("ControlPoint.subscriptionLoop: subscribe successful, will resubscribe in %s intervals", duration)
		return duration
	}
	if err := sub.resubscribe(ctx); err != nil {
		sub.setActiveChan <- false
		duration := 5 * time.Second
		log.Print("ControlPoint.subscriptionLoop:", err)
		return duration
	}
	return sub.getRenewDuration()
}
