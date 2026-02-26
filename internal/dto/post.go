package dto

import (
	"time"
)

type PostResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	ImageURL    *string    `json:"image_url"`
	Abstract    *string    `json:"abstract"`
	Body        string     `json:"body"`
	Published   bool       `json:"published"`
	PublishDate *time.Time `json:"publish_date"`
	Location    *string    `json:"location"`
	Lat         *float64   `json:"lat"`
	Lon         *float64   `json:"lon"`
	Locale      *string    `json:"locale"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	AuthorUsername  string  `json:"author_username"`
	AuthorAvatarURL *string `json:"author_avatar_url"`
	AuthorName      string  `json:"author_name"`

	LikeCount    int  `json:"like_count"`
	UserHasLiked bool `json:"user_has_liked"`
}

type CreatePostRequest struct {
	Title       string     `form:"title" binding:"required"`
	Slug        string     `form:"slug" binding:"required"`
	Abstract    *string    `form:"abstract"`
	Body        string     `form:"body" binding:"required"`
	Published   bool       `form:"published"`
	PublishDate *time.Time `form:"publish_date"`
	Location    *string    `form:"location"`
	Lat         *float64   `form:"lat"`
	Lon         *float64   `form:"lon"`
	Locale      *string    `form:"locale"`
	Tags        []string   `form:"tags"`
}

type UpdatePostRequest struct {
	Title       *string    `form:"title"`
	Slug        *string    `form:"slug"`
	Abstract    *string    `form:"abstract"`
	Body        *string    `form:"body"`
	Published   *bool      `form:"published"`
	PublishDate *time.Time `form:"publish_date"`
	Location    *string    `form:"location"`
	Lat         *float64   `form:"lat"`
	Lon         *float64   `form:"lon"`
	Locale      *string    `form:"locale"`
	Tags        []string   `form:"tags"`
}
