package service

import (
	"context"
	"errors"
	"io"
	"time"

	dbsqlc "san/internal/db/sqlc"
	storage_service "san/internal/service/storage"
	"san/pkg/apperr"
	"san/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type PostService struct {
	repo          dbsqlc.Querier
	activeStorage *storage_service.ActiveStorageService
	log           logger.Logger
}

func NewPostService(repo dbsqlc.Querier, activeStorage *storage_service.ActiveStorageService, log logger.Logger) *PostService {
	return &PostService{
		repo:          repo,
		activeStorage: activeStorage,
		log:           log,
	}
}

type CreatePostInput struct {
	UserID            string
	Title             string
	Slug              string
	ImageFile         io.Reader
	ImageSize         int64
	ImageContentType  string
	ImageOriginalName string
	Abstract          *string
	Body              string
	Published         bool
	PublishDate       *time.Time
	Location          *string
	Lat               *float64
	Lon               *float64
	Locale            *string
	Tags              []string
}

func (s *PostService) CreatePost(ctx context.Context, input CreatePostInput) (*dbsqlc.Post, error) {
	id := uuid.NewString()

	post, err := s.repo.CreatePost(ctx, dbsqlc.CreatePostParams{
		ID:          id,
		UserID:      input.UserID,
		Title:       input.Title,
		Slug:        input.Slug,
		Abstract:    input.Abstract,
		Body:        input.Body,
		Published:   input.Published,
		PublishDate: input.PublishDate,
		Location:    input.Location,
		Lat:         input.Lat,
		Lon:         input.Lon,
		Locale:      input.Locale,
		Tags:        input.Tags,
	})
	if err != nil {
		s.log.Errorf("PostService.CreatePost: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	if input.ImageFile != nil {
		_, err := s.activeStorage.AttachFile(ctx, "posts", post.ID, "image", input.ImageFile, input.ImageSize, input.ImageContentType, input.ImageOriginalName, true)
		if err != nil {
			s.log.Errorf("PostService.CreatePost: failed to upload image: %v", err)
			// Optional: rollback post creation?
			// For now, let's keep the post but return an error or just log it.
			// Ideally rollback.
			if delErr := s.repo.DeletePost(ctx, post.ID); delErr != nil {
				s.log.Errorf("PostService.CreatePost: failed to rollback post creation: %v", delErr)
			}
			return nil, apperr.New(apperr.ErrCodeUploadFailed, "Failed to upload image", 500).WithCause(err)
		}
	}

	return post, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id string) (*dbsqlc.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("Post not found")
		}
		s.log.Errorf("PostService.GetPostByID: %v", err)
		return nil, apperr.InternalServerError(err)
	}
	return post, nil
}

func (s *PostService) ListPosts(ctx context.Context, page, pageSize int32) ([]*dbsqlc.Post, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	posts, err := s.repo.ListPosts(ctx, dbsqlc.ListPostsParams{
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		s.log.Errorf("PostService.ListPosts: %v", err)
		return nil, apperr.InternalServerError(err)
	}
	return posts, nil
}

type UpdatePostInput struct {
	ID                string
	UserID            string // For authorization
	Title             *string
	Slug              *string
	ImageFile         io.Reader
	ImageSize         int64
	ImageContentType  string
	ImageOriginalName string
	Abstract          *string
	Body              *string
	Published         *bool
	PublishDate       *time.Time
	Location          *string
	Lat               *float64
	Lon               *float64
	Locale            *string
	Tags              []string
}

func (s *PostService) UpdatePost(ctx context.Context, input UpdatePostInput) (*dbsqlc.Post, error) {
	// 1. Fetch existing post
	post, err := s.repo.GetPostByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("Post not found")
		}
		s.log.Errorf("PostService.UpdatePost: check post exists: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	// 2. Authorization check
	if post.UserID != input.UserID {
		return nil, apperr.Forbidden("You are not authorized to edit this post")
	}

	// 3. Update
	arg := dbsqlc.UpdatePostParams{
		ID:          input.ID,
		Title:       input.Title,
		Slug:        input.Slug,
		Abstract:    input.Abstract,
		Body:        input.Body,
		Published:   input.Published,
		PublishDate: input.PublishDate,
		Location:    input.Location,
		Lat:         input.Lat,
		Lon:         input.Lon,
		Locale:      input.Locale,
		Tags:        input.Tags,
	}

	updatedPost, err := s.repo.UpdatePost(ctx, arg)
	if err != nil {
		s.log.Errorf("PostService.UpdatePost: %v", err)
		return nil, apperr.InternalServerError(err)
	}

	if input.ImageFile != nil {
		_, err := s.activeStorage.AttachFile(ctx, "posts", post.ID, "image", input.ImageFile, input.ImageSize, input.ImageContentType, input.ImageOriginalName, true)
		if err != nil {
			s.log.Errorf("PostService.UpdatePost: failed to upload image: %v", err)
			return nil, apperr.New(apperr.ErrCodeUploadFailed, "Failed to upload image", 500).WithCause(err)
		}
	}

	return updatedPost, nil
}

func (s *PostService) GetPostImageURL(ctx context.Context, postID string) (string, error) {
	url, err := s.activeStorage.GetAttachmentURL(ctx, "posts", postID, "image")
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		s.log.Debugf("PostService.GetPostImageURL: image not found for post %s: %v", postID, err)
		return "", nil
	}
	return url, nil
}

func (s *PostService) DeletePost(ctx context.Context, id string, userID string) error {
	// 1. Fetch existing post
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperr.NotFound("Post not found")
		}
		s.log.Errorf("PostService.DeletePost: check post exists: %v", err)
		return apperr.InternalServerError(err)
	}

	// 2. Authorization check
	if post.UserID != userID {
		return apperr.Forbidden("You are not authorized to delete this post")
	}

	// 3. Delete
	err = s.repo.DeletePost(ctx, id)
	if err != nil {
		s.log.Errorf("PostService.DeletePost: %v", err)
		return apperr.InternalServerError(err)
	}
	return nil
}
