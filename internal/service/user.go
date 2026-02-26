package service

import (
	"context"
	"errors"
	"io"

	dbsqlc "san/internal/db/sqlc"
	storage_service "san/internal/service/storage"
	"san/pkg/apperr"
	"san/pkg/logger"
	"san/pkg/token"
	"san/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type UserService struct {
	repo          UserRepository
	activeStorage *storage_service.ActiveStorageService
	tokenManager  token.TokenManager
	log           logger.Logger
}

func NewUserService(repo UserRepository, activeStorage *storage_service.ActiveStorageService, tokenManager token.TokenManager, log logger.Logger) *UserService {
	return &UserService{
		repo:          repo,
		activeStorage: activeStorage,
		tokenManager:  tokenManager,
		log:           log,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	User         *dbsqlc.User
	AccessToken  string
	RefreshToken string
}

func (s *UserService) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.Unauthorized("Invalid email or password")
		}
		s.log.Errorf("UserService.Login: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	if err = utils.CheckPassword(input.Password, user.Password); err != nil {
		return nil, apperr.Unauthorized("Invalid email or password")
	}

	accessToken, err := s.tokenManager.CreateAccessToken(user.ID)
	if err != nil {
		s.log.Errorf("UserService.Login: failed to create access token: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	refreshToken, err := s.tokenManager.CreateRefreshToken(user.ID)
	if err != nil {
		s.log.Errorf("UserService.Login: failed to create refresh token: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.tokenManager.VerifyToken(refreshToken)
	if err != nil {
		return "", "", apperr.Unauthorized("Invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return "", "", apperr.Unauthorized("Invalid token type")
	}

	// Verify user exists
	_, err = s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", apperr.Unauthorized("User not found")
		}
		s.log.Errorf("UserService.RefreshToken: %v", err)
		return "", "", apperr.InternalServerError(err)
	}

	newAccessToken, err := s.tokenManager.CreateAccessToken(claims.UserID)
	if err != nil {
		s.log.Errorf("UserService.RefreshToken: failed to create access token: %v", err)
		return "", "", apperr.InternalServerError(err)
	}

	newRefreshToken, err := s.tokenManager.CreateRefreshToken(claims.UserID)
	if err != nil {
		s.log.Errorf("UserService.RefreshToken: failed to create refresh token: %v", err)
		return "", "", apperr.InternalServerError(err)
	}

	return newAccessToken, newRefreshToken, nil
}

type CreateUserInput struct {
	Username           string
	Email              string
	Password           string
	Bio                *string
	AvatarFile         io.Reader
	AvatarSize         int64
	AvatarContentType  string
	AvatarOriginalName string
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*dbsqlc.User, error) {
	id := uuid.NewString()

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, dbsqlc.CreateUserParams{
		ID:       id,
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Bio:      input.Bio,
	})
	if err != nil {
		s.log.Errorf("UserService.CreateUser: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	if input.AvatarFile != nil {
		_, err := s.activeStorage.AttachFile(ctx, "users", user.ID, "avatar", input.AvatarFile, input.AvatarSize, input.AvatarContentType, input.AvatarOriginalName, true)
		if err != nil {
			s.log.Errorf("UserService.CreateUser: failed to upload avatar: %v", err)

			if delErr := s.repo.DeleteUser(ctx, user.ID); delErr != nil {
				s.log.Errorf("UserService.CreateUser: failed to rollback user creation: %v", delErr)
			}

			return nil, apperr.New(apperr.ErrCodeUploadFailed, "Failed to upload avatar", 500).WithCause(err)
		}
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*dbsqlc.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.UserNotFound()
		}
		s.log.Errorf("UserService.GetUserByID: %v", err)
		return nil, apperr.InternalServerError(err)
	}
	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, page, pageSize int32) ([]*dbsqlc.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, err := s.repo.ListUsers(ctx, dbsqlc.ListUsersParams{
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		s.log.Errorf("UserService.ListUsers: %v", err)
		return nil, apperr.InternalServerError(err)
	}
	return users, nil
}

type UpdateUserInput struct {
	ID       string
	Username *string
	Email    *string
	Bio      *string
}

func (s *UserService) UpdateUser(ctx context.Context, input UpdateUserInput) (*dbsqlc.User, error) {
	// Check if user exists
	if _, err := s.repo.GetUserByID(ctx, input.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.UserNotFound()
		}
		s.log.Errorf("UserService.UpdateUser: check user exists: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	arg := dbsqlc.UpdateUserParams{
		ID: input.ID,
	}

	if input.Username != nil {
		arg.Username = input.Username
	}
	if input.Email != nil {
		arg.Email = input.Email
	}
	if input.Bio != nil {
		arg.Bio = input.Bio
	}

	user, err := s.repo.UpdateUser(ctx, arg)
	if err != nil {
		s.log.Errorf("UserService.UpdateUser: %v", err)
		return nil, apperr.InternalServerError(err)
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if _, err := s.repo.GetUserByID(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperr.UserNotFound()
		}
		s.log.Errorf("UserService.DeleteUser: check user exists: %v", err)
		return apperr.InternalServerError(err)
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		s.log.Errorf("UserService.DeleteUser: %v", err)
		return apperr.InternalServerError(err)
	}
	return nil
}

func (s *UserService) GetAvatarURL(ctx context.Context, userID string) (string, error) {
	url, err := s.activeStorage.GetAttachmentURL(ctx, "users", userID, "avatar")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		// Log debug only, as avatar might not exist
		s.log.Debugf("UserService.GetAvatarURL: avatar not found for user %s: %v", userID, err)
		return "", nil
	}
	return url, nil
}

func (s *UserService) UploadUserAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string, originalName string) (*dbsqlc.User, error) {
	_, err := s.activeStorage.AttachFile(ctx, "users", userID, "avatar", file, size, contentType, originalName, true)
	if err != nil {
		s.log.Errorf("UserService.UploadUserAvatar: failed to attach file: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.UserNotFound()
		}
		return nil, apperr.InternalServerError(err)
	}

	return user, nil
}
