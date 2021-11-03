package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func AddStreamRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/streams", func(c *gin.Context) {
		streams, err := p.GetStreams(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, streams)
	})

	r.POST("/stream/new", func(c *gin.Context) {
		// Parse the JSON in the body
		addReq := &api.AddStreamRequest{}
		if err := c.BindJSON(addReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Add the stream
		stream, err := p.AddStream(c, addReq)
		if err != nil {
			if err == api.ErrNameAlreadyExists {
				c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stream)
	})

	r.Use(ensureID)

	r.GET("/stream/:id", func(c *gin.Context) {
		// Get the stream
		stream, err := p.GetStream(c, c.GetInt("id"))
		if err != nil {
			if err == api.ErrStreamNotFound {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stream)
	})

	r.POST("/stream/:id", func(c *gin.Context) {
		// Parse the JSON in the body
		updateReq := &api.UpdateStreamRequest{}
		if err := c.BindJSON(updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		updateReq.ID = c.GetInt("id")

		// Update the stream
		stream, err := p.UpdateStream(c, updateReq)
		if err != nil {
			if err == api.ErrStreamNotFound {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stream)
	})

	r.DELETE("/stream/:id", func(c *gin.Context) {
		// Delete the stream
		err := p.DeleteStream(c, c.GetInt("id"))
		if err != nil {
			if err == api.ErrStreamNotFound {
				c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})
}
