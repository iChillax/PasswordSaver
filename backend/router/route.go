package router

import (
	"backend/engines"

	"github.com/gin-gonic/gin"
)

func CreateRouteTable(app *gin.Engine) {
	v1_group := app.Group("/api/v1")
	route2Health(v1_group)
	route2ManagementBasicAuth(v1_group)
}

func route2Health(group *gin.RouterGroup) {
	group.GET("/healthz", engines.CheckHealth)
}

func route2ManagementBasicAuth(group *gin.RouterGroup) {
	group.POST("/basicauth", engines.CreateBasicAuth)
	group.GET("/basicauth", engines.ListAllBasicAuth)
}
