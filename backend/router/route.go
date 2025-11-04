package router

import (
	"backend/engines"

	"github.com/gin-gonic/gin"
)

func CreateRouteTable(app *gin.Engine) {
	app.GET("/healthz", engines.CheckHealth)
}
