package routes

import (
	"url-shortener/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/shorten", controllers.ShortenURL)
	r.GET("/:shortCode", controllers.RedirectURL)

	return r
}
