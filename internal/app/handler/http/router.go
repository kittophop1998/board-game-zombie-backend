package http

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	apiV1 := router.Group("/api/v1")
	{
		userRoutes := apiV1.Group("/users")
		{
			userRoutes.GET("", H.User.GetUsers)
		}
	}
}
