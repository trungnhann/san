package service

import (
	"context"
	"fmt"
	"io"

	dbsqlc "san/internal/db/sqlc"
	storage_service "san/internal/service/storage"
	"san/pkg/logger"
	"san/pkg/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo          UserRepository
	activeStorage *storage_service.ActiveStorageService
	log           logger.Logger
}

func NewUserService(repo UserRepository, activeStorage *storage_service.ActiveStorageService, log logger.Logger) *UserService {
	return &UserService{
		repo:          repo,
		activeStorage: activeStorage,
		log:           log,
	}
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
		return nil, err
	}

	if input.AvatarFile != nil {
		_, err := s.activeStorage.AttachFile(ctx, "users", user.ID, "avatar", input.AvatarFile, input.AvatarSize, input.AvatarContentType, input.AvatarOriginalName, true)
		if err != nil {
			s.log.Errorf("UserService.CreateUser: failed to upload avatar: %v", err)

			if delErr := s.repo.DeleteUser(ctx, user.ID); delErr != nil {
				s.log.Errorf("UserService.CreateUser: failed to rollback user creation: %v", delErr)
			}

			return nil, fmt.Errorf("failed to upload avatar: %w", err)
		}
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*dbsqlc.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.log.Errorf("UserService.GetUserByID: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*dbsqlc.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		s.log.Errorf("UserService.ListUsers: %v", err)
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetAvatarURL(ctx context.Context, userID string) (string, error) {
	url, err := s.activeStorage.GetAttachmentURL(ctx, "users", userID, "avatar")
	if err != nil {
		s.log.Debugf("UserService.GetAvatarURL: avatar not found for user %s: %v", userID, err)
		return "", nil
	}
	return url, nil
}

func (s *UserService) UploadUserAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string, originalName string) (*dbsqlc.User, error) {
	_, err := s.activeStorage.AttachFile(ctx, "users", userID, "avatar", file, size, contentType, originalName, true)
	if err != nil {
		s.log.Errorf("UserService.UploadUserAvatar: failed to attach file: %v", err)
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
