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

func AddRadioRoutes(r *gin.RouterGroup, a *api.API, upgrader *websocket.Upgrader) {
	r.GET("/radios", func(c *gin.Context) {
		c.JSON(http.StatusOK, a.GetRadioStates())
	})

	r.POST("/radios", func(c *gin.Context) {
		err := a.DiscoverRadios()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"msg": err.Error()})
			return
		}
	})

	r.GET("/radio/ws", func(c *gin.Context) {
		// Upgrade connection to websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}

		// Handle websocket
		a.HandleWS(conn, "")
	})

	r.Use(ensureUUID(a))

	r.GET("/radio/:UUID", func(c *gin.Context) {
		// Get UUID
		uuid, _ := c.Params.Get("UUID")

		// Get Radio or return 404
		state, ok := a.GetRadioState(uuid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Return Radio
		c.JSON(http.StatusOK, state)
	})

	r.PATCH("/radio/:UUID", func(c *gin.Context) {
		// Get UUID
		uuid, _ := c.Params.Get("UUID")

		// Get Radio or return 404
		rd, ok := a.GetRadio(uuid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Parse json body
		var radioPost RadioPost
		err := c.BindJSON(&radioPost)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// TODO: Add shorter ctx timeout, currently soap timeout is 30 seconds

		// Set power if not nil
		if radioPost.Power != nil {
			if err := rd.SetPowerState(c, *radioPost.Power); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err})
				return
			}
		}

		// Play preset if not nil
		if radioPost.Preset != nil {
			if !rd.IsPresetValid(*radioPost.Preset) {
				c.JSON(http.StatusBadRequest, gin.H{"err": "preset is not valid"})
				return
			}
			if err := rd.PlayPreset(c, *radioPost.Preset); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err})
				return
			}
		}

		// Set volume if not nil
		if radioPost.Volume != nil {
			if err := rd.SetVolume(c, radio.NormalizeVolume(*radioPost.Volume)); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err})
				return
			}
		}

		c.JSON(http.StatusOK, radioPost)
	})

	r.GET("/radio/:UUID/ws", func(c *gin.Context) {
		// Get UUID
		uuid, _ := c.Params.Get("UUID")

		// Return 404 if radio does not exist
		if !a.IsValidRadio(uuid) {
			c.Status(http.StatusNotFound)
			return
		}

		// Upgrade connection to websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}

		// Handle websocket
		a.HandleWS(conn, uuid)
	})
}
