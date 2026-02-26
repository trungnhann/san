package router

import (
	"san/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("/:id/avatar", userHandler.UploadAvatar)
	}
}
