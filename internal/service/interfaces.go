package service

import (
	"context"

	dbsqlc "san/internal/db/sqlc"
)

// UserRepository defines the interface for user data access
// In Clean Architecture, the Service layer defines the interface it needs (Dependency Inversion)
// Ideally this would match dbsqlc.Querier, or be a subset of it.
type UserRepository interface {
	CreateUser(ctx context.Context, arg dbsqlc.CreateUserParams) (*dbsqlc.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserByID(ctx context.Context, id string) (*dbsqlc.User, error)
	ListUsers(ctx context.Context) ([]*dbsqlc.User, error)
}
