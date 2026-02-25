package router

import (
	"san/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	api := r.Group("/api")
	{
		users := api.Group("/users")
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("/:id/avatar", userHandler.UploadAvatar)
	}
}
