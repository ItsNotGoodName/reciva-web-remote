package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/go-chi/chi/v5"
)

type RadioRequester func(*http.Request, radio.Radio) presenter.Response

func handleRadioError(err error) presenter.Response {
	code := http.StatusInternalServerError
	if err == radio.ErrRadioClosed {
		code = http.StatusGone
	} else if err == radio.ErrRadioNotFound {
		code = http.StatusNotFound
	}

	return presenter.Response{
		Code:  code,
		Error: err,
	}
}

func RequireRadio(hub radio.HubService, next RadioRequester) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		uuid := chi.URLParam(r, "uuid")

		rd, err := hub.Get(uuid)
		if err != nil {
			return handleRadioError(err)
		}

		return next(r, rd)
	}
}

func GetRadios(hub radio.HubService, radioService radio.RadioService) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		// List radios
		radios := hub.List()

		// List states
		states := make([]state.State, 0, len(radios))
		for _, rd := range radios {
			state, err := radioService.GetState(r.Context(), rd)
			if err != nil {
				log.Println("api.GetRadios:", err)
				continue
			}
			states = append(states, *state)
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: states,
		}
	}
}

func PostRadios(hub radio.HubService) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		// Discover radios
		count, err := hub.Discover()
		if err != nil {
			code := http.StatusInternalServerError
			if err == radio.ErrHubDiscovering {
				code = http.StatusConflict
			} else if err == radio.ErrHubServiceClosed {
				code = http.StatusServiceUnavailable
			}
			return presenter.Response{
				Code:  code,
				Error: err,
			}
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: count,
		}
	}
}

func GetRadio(radioService radio.RadioService) RadioRequester {
	return func(r *http.Request, rd radio.Radio) presenter.Response {
		// Get state
		state, err := radioService.GetState(r.Context(), rd)
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: state,
		}
	}
}

func PatchRadio(radioService radio.RadioService) RadioRequester {
	type RadioPatch struct {
		Power  *bool `json:"power,omitempty"`
		Preset *int  `json:"preset,omitempty"`
		Volume *int  `json:"volume,omitempty"`
	}

	return func(r *http.Request, rd radio.Radio) presenter.Response {
		// Parse body
		var radioPatch RadioPatch
		if err := json.NewDecoder(r.Body).Decode(&radioPatch); err != nil {
			return presenter.Response{
				Code:  http.StatusBadRequest,
				Error: err,
			}
		}

		// Set radio power
		if radioPatch.Power != nil {
			if err := radioService.SetPower(r.Context(), rd, *radioPatch.Power); err != nil {
				return handleRadioError(err)
			}
		}
		// Set radio preset
		if radioPatch.Preset != nil {
			if err := radioService.PlayPreset(r.Context(), rd, *radioPatch.Preset); err != nil {
				return handleRadioError(err)
			}
		}
		// Set radio volume
		if radioPatch.Volume != nil {
			if err := radioService.SetVolume(r.Context(), rd, *radioPatch.Volume); err != nil {
				return handleRadioError(err)
			}
		}

		return presenter.Response{
			Code: http.StatusOK,
		}
	}
}

func PostRadio(radioService radio.RadioService) RadioRequester {
	return func(r *http.Request, rd radio.Radio) presenter.Response {
		// Refresh radio
		if err := radioService.Refresh(r.Context(), rd); err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
		}
	}
}

func PostRadioVolume(radioService radio.RadioService) RadioRequester {
	return func(r *http.Request, rd radio.Radio) presenter.Response {
		// Refresh radio volume
		if err := radioService.RefreshVolume(r.Context(), rd); err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
		}
	}
}
