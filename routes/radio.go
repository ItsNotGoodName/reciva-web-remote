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
		c.JSON(http.StatusOK, a.GetRadioStates(c))
	})

	r.POST("/radios", func(c *gin.Context) {
		err := a.DiscoverRadios()
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"msg": err.Error()})
			return
		}
	})

	r.GET("/radio/ws", func(c *gin.Context) {
		// Get UUID
		uuid, ok := c.GetQuery("uuid")
		if ok {
			// Return 404 if radio does not exist
			if !a.IsValidRadio(uuid) {
				c.Status(http.StatusNotFound)
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
		a.HandleWS(conn, uuid)
	})

	r.Use(ensureUUID(a))

	r.GET("/radio/:UUID", func(c *gin.Context) {
		// Get UUID
		uuid, _ := c.Params.Get("UUID")

		// Get Radio or return 404
		state, ok := a.GetRadioState(c, uuid)
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

		// Set power if not nil
		if radioPost.Power != nil {
			if err := rd.SetPower(c, *radioPost.Power); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err})
				return
			}
		}

		// Play preset if not nil
		if radioPost.Preset != nil {
			if err := rd.PlayPreset(c, *radioPost.Preset); err != nil {
				if err == radio.ErrInvalidPreset {
					c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
					return
				}
				c.JSON(http.StatusServiceUnavailable, gin.H{"err": err.Error()})
				return
			}
		}

		// Set volume if not nil
		if radioPost.Volume != nil {
			if err := rd.SetVolume(*radioPost.Volume); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, radioPost)
	})

	r.POST("/radio/:UUID/renew", func(c *gin.Context) {
		// Get UUID
		uuid, _ := c.Params.Get("UUID")

		// Return 404 if radio does not exist
		rd, ok := a.GetRadio(uuid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Renew
		rd.Subscription.Renew()
	})

	r.POST("/radio/:UUID/volume", func(c *gin.Context) {
		// Get UUID
		uuid, _ := c.Params.Get("UUID")

		// Return 404 if radio does not exist
		rd, ok := a.GetRadio(uuid)
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		// Refresh volume
		if err := rd.RefreshVolume(c); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
	})
}
