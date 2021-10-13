package goupnpsub

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Renew tells subscription to renew if it is not already renewing.
func (sub *Subscription) Renew() {
	select {
	case sub.renewChan <- true:
	default:
		return
	}
}

// unsubscribe sends an UNSUBSCRIBE request to event publisher.
func (sub *Subscription) unsubscribe(ctx context.Context) error {
	// Create request
	req, err := http.NewRequest("UNSUBSCRIBE", sub.eventUrl, nil)
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

	// Check if request failed
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("subscribe: response's status code is invalid, %s", res.Status)
	}

	return nil
}

// resubscribe sends a SUBSCRIBE request to event publisher that renews the existing subscription.
func (sub *Subscription) resubscribe(ctx context.Context) error {
	// Create request
	req, err := http.NewRequest("SUBSCRIBE", sub.eventUrl, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	// Add headers to request
	req.Header.Add("SID", sub.sid)
	req.Header.Add("TIMEOUT", defaultTimeout)

	// Execute request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check if request failed
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("resubscribe: response's status code is invalid, %s", res.Status)
	}

	// Check request's SID header
	resSid := res.Header.Get("SID")
	if resSid == "" {
		return errors.New("resubscribe: response did not supply a sid")
	}
	if resSid != sub.sid {
		return fmt.Errorf("resubscribe: response's sid does not match sub's sid, %s != %s", resSid, sub.sid)
	}

	// Update sub's timeout with request's SID
	timeout, err := parseTimeout(res.Header.Get("timeout"))
	if err != nil {
		return err
	}
	sub.timeout = timeout

	return nil
}
