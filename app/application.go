package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

// StartApplication ..
func StartApplication() {
	mapUrls()
	router.NoRoute(func(c *gin.Context) { c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"}) })
	router.Run(":8080")
}
