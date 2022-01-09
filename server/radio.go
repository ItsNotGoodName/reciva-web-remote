package server

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleGetRadio() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		// Get state
		state, err := rd.GetState(c)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
			return
		}

		// Return Radio
		c.JSON(http.StatusOK, state)
	}
}

func (s *Server) handleGetRadios() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, s.h.GetRadioStates(c))
	}
}

func (s *Server) handleGetRadioWS() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get uuid
		uuid, ok := c.GetQuery("uuid")
		if ok {
			// Return 404 if radio does not exist
			if !s.h.IsValidRadio(uuid) {
				c.JSON(http.StatusNotFound, gin.H{"err": radio.ErrRadioNotFound.Error()})
				return
			}
		}

		// Upgrade connection to websocket
		conn, err := s.u.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		// Handle websocket
		api.NewRadioWS(conn, s.h).Start(uuid)
	}
}

func (s *Server) handlePatchRadio() func(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Set power if not nil and preset is nil
		if radioPatch.Preset == nil {
			if radioPatch.Power != nil {
				if err := rd.SetPower(c, *radioPatch.Power); err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
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
				c.JSON(code, gin.H{"err": err.Error()})
				return
			}
		}

		// Set volume if not nil
		if radioPatch.Volume != nil {
			if err := rd.SetVolume(c, *radioPatch.Volume); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
				return
			}
		}
	}
}

func (s *Server) handlePostRadio() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		rd.Refresh()
	}
}

func (s *Server) handlePostRadioVolume() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get radio
		rd := c.MustGet("radio").(*radio.Radio)

		// Refresh volume
		if err := rd.RefreshVolume(c); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
			return
		}
	}
}

func (s *Server) handlePostRadios() func(c *gin.Context) {
	return func(c *gin.Context) {
		err := s.h.Discover()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
	}
}
