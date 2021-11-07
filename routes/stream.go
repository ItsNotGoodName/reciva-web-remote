package routes

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/gin-gonic/gin"
)

func AddStreamRoutes(r *gin.RouterGroup, p *api.PresetAPI) {
	r.GET("/streams", func(c *gin.Context) {
		// Read all streams
		streams, err := p.ReadStreams(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, streams)
	})

	r.POST("/stream/new", func(c *gin.Context) {
		// Parse the JSON in the body
		createReq := &api.CreateStreamRequest{}
		if err := c.BindJSON(createReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		// Create the stream
		stream, err := p.CreateStream(c, createReq)
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrStreamNameInvalid || err == api.ErrStreamContentInvalid {
				code = http.StatusBadRequest
			}
			if err == api.ErrStreamNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stream)
	})

	r.Use(ensureID)

	r.GET("/stream/:id", func(c *gin.Context) {
		// Read the stream
		stream, err := p.ReadStream(c, c.GetInt("id"))
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrStreamNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
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
			code := http.StatusInternalServerError
			if err == api.ErrStreamNameInvalid || err == api.ErrStreamContentInvalid {
				code = http.StatusBadRequest
			}
			if err == api.ErrStreamNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stream)
	})

	r.DELETE("/stream/:id", func(c *gin.Context) {
		// Delete the stream
		err := p.DeleteStream(c, c.GetInt("id"))
		if err != nil {
			code := http.StatusInternalServerError
			if err == api.ErrStreamNotFound {
				code = http.StatusNotFound
			}
			c.JSON(code, gin.H{"err": err.Error()})
			return
		}
	})
}
