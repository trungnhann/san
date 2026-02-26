package dto

import "time"

type CreateUserRequest struct {
	Username string  `json:"username" form:"username" binding:"required"`
	Email    string  `json:"email" form:"email" binding:"required,email"`
	Password string  `json:"password" form:"password" binding:"required"`
	Bio      *string `json:"bio" form:"bio"`
	Image    *string `json:"image" form:"image"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Bio       *string   `json:"bio"`
	Image     *string   `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
