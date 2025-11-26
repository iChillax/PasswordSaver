package main

import (
	"backend/middleware"
	"backend/router"
	"backend/settings"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	// Add CORS middleware
	app.Use(middleware.CORSMiddleware())

	settings.Initiate()
	router.CreateRouteTable(app)
	app.Run("0.0.0.0:8080")
}
