package radio

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
)

func ListRadios(h *hub.Hub) []model.Radio {
	hubRadios := h.List()
	radios := make([]model.Radio, len(hubRadios))
	for i, r := range hubRadios {
		radios[i] = model.NewRadio(r)
	}

	return radios
}

func ListStates(ctx context.Context, h *hub.Hub) []state.State {
	rds := h.List()
	var states []state.State
	for _, p := range rds {
		s, err := p.State(ctx)
		if err != nil {
			log.Println("radio.ListStates:", err)
			continue
		}

		states = append(states, *s)
	}

	return states
}

func DeleteRadio(h *hub.Hub, r hub.Radio) error {
	if err := h.Delete(r.UUID); err != nil {
		return err
	}

	pubsub.PublishStaleRadios(pubsub.DefaultPub)

	return nil
}
