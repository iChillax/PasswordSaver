package main

import (
	"backend/router"
	"backend/settings"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	settings.Initiate()
	router.CreateRouteTable(app)
	app.Run("0.0.0.0:8000")
}
