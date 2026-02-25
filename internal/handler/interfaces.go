package handler

import (
	"context"
	"io"

	dbsqlc "san/internal/db/sqlc"
	"san/internal/service"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, input service.CreateUserInput) (*dbsqlc.User, error)
	GetUserByID(ctx context.Context, id string) (*dbsqlc.User, error)
	ListUsers(ctx context.Context) ([]*dbsqlc.User, error)
	UploadUserAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string, originalName string) (*dbsqlc.User, error)
	GetAvatarURL(ctx context.Context, userID string) (string, error)
}
