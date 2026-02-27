package router

import (
	"san/internal/handler"
	"san/internal/middleware"
	"san/pkg/token"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *handler.UserHandler, postHandler *handler.PostHandler, tokenManager token.TokenManager) {
	v1 := r.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)

		users := v1.Group("/users")
		users.POST("", userHandler.CreateUser)
		users.POST("/verify", userHandler.VerifyEmail)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUserByID)

		posts := v1.Group("/posts")
		posts.GET("", postHandler.ListPosts)
		posts.GET("/:id", postHandler.GetPostByID)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(tokenManager))
		{
			protectedUsers := protected.Group("/users")
			protectedUsers.PUT("/:id", userHandler.UpdateUser)
			protectedUsers.DELETE("/:id", userHandler.DeleteUser)
			protectedUsers.POST("/:id/avatar", userHandler.UploadAvatar)

			protectedPosts := protected.Group("/posts")
			protectedPosts.POST("", postHandler.CreatePost)
			protectedPosts.PUT("/:id", postHandler.UpdatePost)
			protectedPosts.DELETE("/:id", postHandler.DeletePost)
		}
	}
}
