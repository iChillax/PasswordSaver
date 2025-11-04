package main

import (
	"backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	router.CreateRouteTable(app)
	app.Run("0.0.0.0:8000")
}
