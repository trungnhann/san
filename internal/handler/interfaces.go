package handler

import (
	"context"
	"io"

	dbsqlc "san/internal/db/sqlc"
	"san/internal/service"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	CreateUser(ctx context.Context, input service.CreateUserInput) (*dbsqlc.User, error)
	GetUserByID(ctx context.Context, id string) (*dbsqlc.User, error)
	ListUsers(ctx context.Context, page, pageSize int32) ([]*dbsqlc.User, error)
	UploadUserAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string, originalName string) (*dbsqlc.User, error)
	GetAvatarURL(ctx context.Context, userID string) (string, error)
	UpdateUser(ctx context.Context, input service.UpdateUserInput) (*dbsqlc.User, error)
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, input service.LoginInput) (*service.LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	VerifyEmail(ctx context.Context, email string, otp string) error
}
