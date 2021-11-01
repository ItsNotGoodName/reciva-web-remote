package goupnpsub

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Renew tells subscription to renew if it is not already renewing.
func (sub *Subscription) Renew() {
	select {
	case sub.renewChan <- true:
	default:
	}
}

// activeLoop handles active status of subscription.
func (sub *Subscription) activeLoop() {
	log.Println("Subscription.activeLoop: started")

	active := false
	for {
		select {
		case <-sub.Done:
			close(sub.Active)
			return
		case active = <-sub.setActiveChan:
		case sub.Active <- active:
		}
	}
}

// unsubscribe sends an UNSUBSCRIBE request to event publisher.
func (sub *Subscription) unsubscribe(ctx context.Context) error {
	// Create request
	req, err := http.NewRequest("UNSUBSCRIBE", sub.eventURL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// Add headers to request
	req.Header.Add("SID", sub.sid)

	// Execute request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check if request failed
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("subscribe: response's status code is invalid, %s", res.Status)
	}

	return nil
}

// resubscribe sends a SUBSCRIBE request to event publisher that renews the existing subscription.
func (sub *Subscription) resubscribe(ctx context.Context) error {
	// Create request
	req, err := http.NewRequest("SUBSCRIBE", sub.eventURL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// Add headers to request
	req.Header.Add("SID", sub.sid)
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
		return fmt.Errorf("resubscribe: response's status code is invalid, %s", res.Status)
	}

	// Check request's SID header
	sid := res.Header.Get("SID")
	if sid == "" {
		return errors.New("resubscribe: response did not supply a sid")
	}
	if sid != sub.sid {
		return fmt.Errorf("resubscribe: response's sid does not match sub's sid, %s != %s", sid, sub.sid)
	}

	// Update sub's timeout
	timeout, err := unmarshalTimeout(res.Header.Get("timeout"))
	if err != nil {
		return err
	}
	sub.timeout = timeout

	return nil
}

// setActive sets active status of subscription.
func (sub *Subscription) setActive(ctx context.Context, active bool) {
	select {
	case <-ctx.Done():
	case sub.setActiveChan <- active:
	}
}
