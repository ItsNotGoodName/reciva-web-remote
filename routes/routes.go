package routes

import (
	"log"
	"strconv"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/gin-gonic/gin"
)

func Start(cfg *config.Config, r *gin.Engine) {
	log.Println("routes.Start: listening on port", cfg.Port)
	log.Fatal("routes.Start", r.Run(":"+strconv.Itoa(cfg.Port)))
}
