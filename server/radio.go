package server

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func handleRadioGet() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		// Get state
		state, err := rd.GetState(c)
		if err != nil {
			renderError(c, http.StatusServiceUnavailable, err)
			return
		}

		// Return Radio
		c.JSON(http.StatusOK, state)
	}
}

func handleRadioList(h *radio.Hub) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, h.GetRadioStates(c))
	}
}

func handleRadioWS(hub *radio.Hub, upgrader *websocket.Upgrader) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get uuid
		uuid, _ := c.GetQuery("uuid")

		// Upgrade connection to websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			renderError(c, http.StatusInternalServerError, err)
			return
		}

		// Handle websocket
		go ws.Handle(conn, hub, uuid)
	}
}

func handleRadioPatch() func(c *gin.Context) {
	type RadioPatch struct {
		Power  *bool `json:"power,omitempty"`
		Preset *int  `json:"preset,omitempty"`
		Volume *int  `json:"volume,omitempty"`
	}

	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		// Parse json body
		var radioPatch RadioPatch
		err := c.BindJSON(&radioPatch)
		if err != nil {
			renderError(c, http.StatusBadRequest, err)
			return
		}

		// Set power if not nil and preset is nil
		if radioPatch.Preset == nil {
			if radioPatch.Power != nil {
				if err := rd.SetPower(c, *radioPatch.Power); err != nil {
					renderError(c, http.StatusServiceUnavailable, err)
					return
				}
			}
		} else {
			// Play preset if not nil
			if err := rd.PlayPreset(c, *radioPatch.Preset); err != nil {
				code := http.StatusServiceUnavailable
				if err == radio.ErrPresetInvalid {
					code = http.StatusBadRequest
				}
				renderError(c, code, err)
				return
			}
		}

		// Set volume if not nil
		if radioPatch.Volume != nil {
			if err := rd.SetVolume(c, *radioPatch.Volume); err != nil {
				renderError(c, http.StatusServiceUnavailable, err)
				return
			}
		}
	}
}

func handleRadioRefresh() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		rd.Refresh()
	}
}

func handleRadioVolumeRefresh() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		// Refresh volume
		if err := rd.RefreshVolume(c); err != nil {
			renderError(c, http.StatusServiceUnavailable, err)
			return
		}
	}
}

func handleRadioDiscover(h *radio.Hub) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := h.Discover(); err != nil {
			renderError(c, http.StatusConflict, err)
			return
		}
	}
}
