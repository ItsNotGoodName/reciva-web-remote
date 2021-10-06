package api

import (
	"log"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gorilla/websocket"
)

func (a *API) HandleRadioWS(conn *websocket.Conn, uuid string) {
	c := radio.HubClient{Send: make(chan radio.State)}
	a.h.Register <- &c
	go a.writeRadioWS(&c, conn, uuid)
	go a.readRadioWS(&c, conn, uuid)
}

func (a *API) writeRadioWS(c *radio.HubClient, conn *websocket.Conn, uuid string) {
	defer func() {
		a.h.Unregister <- c
		conn.Close()
	}()

	state, ok := a.GetRadioState(uuid)
	if !ok {
		log.Printf("writeRadioWS: could not find state with uuid, %s", uuid)
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		return
	}

	if err := conn.WriteJSON(state); err != nil {
		log.Println(err)
		return
	}

	for state := range c.Send {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if state.UUID != uuid {
			continue
		}

		if err := conn.WriteJSON(state); err != nil {
			log.Println(err)
			return
		}
	}
	conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (a *API) readRadioWS(c *radio.HubClient, conn *websocket.Conn, uuid string) {
	defer func() {
		a.h.Unregister <- c
		conn.Close()
	}()

	conn.SetReadLimit(512)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			return
		}
	}
}
