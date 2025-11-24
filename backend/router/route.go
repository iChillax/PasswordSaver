package router

import (
	"backend/engines"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func CreateRouteTable(app *gin.Engine) {
	v1_group := app.Group("/api/v1")
	route2Health(v1_group)
	route2ManagementBasicAuth(v1_group)
	route2Auth(v1_group)
	route2Secrets(v1_group)
}

func route2Health(group *gin.RouterGroup) {
	group.GET("/healthz", engines.CheckHealth)
}

func route2ManagementBasicAuth(group *gin.RouterGroup) {
	group.POST("/basicauth", engines.CreateBasicAuth)
	group.GET("/basicauth", engines.ListAllBasicAuth)
}

func route2Auth(group *gin.RouterGroup) {
	group.POST("/auth/register", engines.Register)
	group.POST("/auth/login", engines.Login)
}

func route2Secrets(group *gin.RouterGroup) {
	// Apply auth middleware to all secret routes
	secretsGroup := group.Group("/secrets")
	secretsGroup.Use(middleware.AuthMiddleware())
	{
		secretsGroup.POST("", engines.CreateSecret)
		secretsGroup.GET("", engines.ListSecrets)
		secretsGroup.GET("/search", engines.SearchSecrets)
		secretsGroup.GET("/:id", engines.GetSecret)
		secretsGroup.PUT("/:id", engines.UpdateSecret)
		secretsGroup.DELETE("/:id", engines.DeleteSecret)
	}
}
