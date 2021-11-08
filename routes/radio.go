package routes

import (
	"log"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RadioPost struct {
	Power  *bool `json:"power,omitempty"`
	Preset *int  `json:"preset,omitempty"`
	Volume *int  `json:"volume,omitempty"`
}

func AddRadioRoutes(r *gin.RouterGroup, h *radio.Hub, upgrader *websocket.Upgrader) {
	r.GET("/radios", func(c *gin.Context) {
		c.JSON(http.StatusOK, h.GetRadioStates(c))
	})

	r.POST("/radios", func(c *gin.Context) {
		err := h.Discover()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
			return
		}
	})

	r.GET("/radio/ws", func(c *gin.Context) {
		// Get uuid
		uuid, ok := c.GetQuery("uuid")
		if ok {
			// Return 404 if radio does not exist
			if !h.IsValidRadio(uuid) {
				c.JSON(http.StatusNotFound, gin.H{"err": radio.ErrRadioNotFound.Error()})
				return
			}
		}

		// Upgrade connection to websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Routes.AddRadioRoutes(ERROR):", err)
			return
		}

		// Handle websocket
		api.NewRadioWS(conn, h).Start(uuid)
	})

	r.Use(ensureUUID)

	r.GET("/radio/:uuid", func(c *gin.Context) {
		// Get uuid
		uuid := c.GetString("uuid")

		// Get Radio or return 404
		state, err := h.GetRadioState(c, uuid)
		if err != nil {
			code := http.StatusInternalServerError
			if err == radio.ErrRadioNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		// Return Radio
		c.JSON(http.StatusOK, state)
	})

	r.PATCH("/radio/:uuid", func(c *gin.Context) {
		// Get uuid
		uuid := c.GetString("uuid")

		// Get Radio or return 404
		rd, ok := h.GetRadio(uuid)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"err": radio.ErrRadioNotFound.Error()})
			return
		}

		// Parse json body
		var radioPost RadioPost
		err := c.BindJSON(&radioPost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Set power if not nil and preset is nil
		if radioPost.Preset == nil {
			if radioPost.Power != nil {
				if err := rd.SetPower(c, *radioPost.Power); err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
					return
				}
			}
		} else {
			// Play preset if not nil
			if err := rd.PlayPreset(c, *radioPost.Preset); err != nil {
				code := http.StatusServiceUnavailable
				if err == radio.ErrInvalidPreset {
					code = http.StatusBadRequest
				}
				c.JSON(code, gin.H{"err": err.Error()})
				return
			}
		}

		// Set volume if not nil
		if radioPost.Volume != nil {
			if err := rd.SetVolume(*radioPost.Volume); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
				return
			}
		}
	})

	r.POST("/radio/:uuid/renew", func(c *gin.Context) {
		// Get uuid
		uuid := c.GetString("uuid")

		// Return 404 if radio does not exist
		rd, ok := h.GetRadio(uuid)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"err": radio.ErrRadioNotFound.Error()})
			return
		}

		// Renew
		rd.Renew()
	})

	r.POST("/radio/:uuid/volume", func(c *gin.Context) {
		// Get uuid
		uuid := c.GetString("uuid")

		// Return 404 if radio does not exist
		rd, ok := h.GetRadio(uuid)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"err": radio.ErrRadioNotFound.Error()})
			return
		}

		// Refresh volume
		if err := rd.RefreshVolume(c); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
			return
		}
	})
}
