package goupnpsub

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// NewControlPoint creates and returns a ControlPoint.
func NewControlPoint() *ControlPoint {
	return NewControlPointWithPort("8058")
}

// NewControlPoint creates and returns a ControlPoint that listens on a specific port.
func NewControlPointWithPort(listenPort string) *ControlPoint {
	cp := &ControlPoint{
		listenUri:     "/eventSub",
		listenPort:    listenPort,
		sidMap:        make(map[string]*Subscription),
		sidMapRWMutex: sync.RWMutex{},
	}
	http.Handle(cp.listenUri, cp)
	return cp
}

// Start creates a http server that listens for notify requests.
func (cp *ControlPoint) Start() {
	log.Fatal(http.ListenAndServe(":"+cp.listenPort, nil))
}

// NewSubscription creates and returns a Subscription.
func (cp *ControlPoint) NewSubscription(dctx context.Context, eventUrl *url.URL) (*Subscription, error) {
	// Get callback url
	callbackUrl, err := getCallbackIP(eventUrl)
	if err != nil {
		return nil, err
	}

	// Crete sub and start subscriptionLoop
	sub := &Subscription{
		ActiveChan:  make(chan bool),
		EventChan:   make(chan *Event, 10),
		renewChan:   make(chan bool),
		eventUrl:    eventUrl.String(),
		callbackUrl: "<http://" + callbackUrl + ":" + cp.listenPort + cp.listenUri + ">",
	}
	go cp.subscriptionLoop(dctx, sub)

	return sub, nil
}

// ServeHTTP handles notify requests from event publishers and is called by the http module.
func (cp *ControlPoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Do I need w.Write(nil)?
	defer w.Write(nil)

	// TODO: Validate request

	// Get SID from request
	sid := r.Header.Get("SID")
	if sid == "" {
		log.Println("ServeHTTP: notify request did not supply sid")
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// Get body from from request
	// TODO: Do I need r.Body.Close()?
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Parse xmlEvent from request's body
	xmlEvent, err := parseEventXML(body)
	if err != nil {
		log.Println(err)
		return
	}

	// Parse properties from xmlEvent
	properties := parseProperties(xmlEvent)

	// Create event using SID and properties
	event := Event{sid: sid, Properties: properties}

	// TODO: Move this section up if it does not causes locking issues with ControlPoint.subscribe() function
	// Find sub from sidMap using SID
	cp.sidMapRWMutex.RLock()
	sub, ok := cp.sidMap[sid]
	cp.sidMapRWMutex.RUnlock()
	if !ok {
		log.Println("ServeHTTP(WARNING): could not find sid in sidMap")
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// Try to send event to sub's EventChan, fail after waiting for 20 seconds
	select {
	case <-time.After(20 * time.Second):
		log.Println("ServeHTTP(ERROR): could not send event to subscription's EventChan")
	case sub.EventChan <- &event:
	}
}

// subscribe sends SUBSCRIBE request to event publisher to create a new subscription.
func (cp *ControlPoint) subscribe(ctx context.Context, sub *Subscription) error {
	// Create request
	req, err := http.NewRequest("SUBSCRIBE", sub.eventUrl, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// Add headers to request
	req.Header.Add("CALLBACK", sub.callbackUrl)
	req.Header.Add("NT", "upnp:event")
	req.Header.Add("TIMEOUT", DefaultTimeout)

	// Execute request
	client := http.Client{}
	res, err := client.Do(req)

	// Lock map as soon as possible to prevent a race condition with upnp publisher
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
func (cp *ControlPoint) subscriptionLoop(dctx context.Context, sub *Subscription) {
	// TODO: Refactor this function
	log.Println("subscriptionLoop: started")

	// Subscription status goroutine
	activeChan := make(chan bool)
	go func() {
		active := false
		for {
			select {
			case <-dctx.Done():
				return
			case active = <-activeChan:
			case sub.ActiveChan <- active:
			}
		}
	}()

	// Renew
	var duration time.Duration
	renew := func() {
		if !<-sub.ActiveChan {
			if err := cp.subscribe(dctx, sub); err != nil {
				log.Print(err)
				return
			}
			activeChan <- true
			duration = getRenewDuration(sub.timeout)
			log.Printf("subscriptionLoop: subscribe successful, will resubscribe in %s", duration)
			return
		}
		if err := sub.resubscribe(dctx); err != nil {
			activeChan <- false
			duration = 5 * time.Second
			log.Println(err)
			return
		}
		duration = getRenewDuration(sub.timeout)
		log.Printf("subscriptionLoop: resubscribe successful, will resubscribe in %s", duration)
	}
	renew()

	for {
		select {
		case <-dctx.Done():
			log.Println("subscriptionLoop: dctx is done, starting cleanup")

			// Delete sub.sid from sidMap
			cp.sidMapRWMutex.Lock()
			delete(cp.sidMap, sub.sid)
			cp.sidMapRWMutex.Unlock()

			// Unsubscribe
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			err := sub.unsubscribe(ctx)
			if err != nil {
				log.Println(err)
			}

			log.Println("subscriptionLoop: cleanup finished")
			return
		case <-sub.renewChan:
			log.Println("subscriptionLoop: RenewChan received")
			renew()
		case <-time.After(duration): // TODO: Use time.NewTimer to prevent memory usage when RenewChan spam called
			renew()
		}
	}
}
