package handler

import (
	"context"
	"net/http"
	"strconv"

	dbsqlc "san/internal/db/sqlc"
	"san/internal/dto"
	"san/internal/service"
	"san/pkg/apperr"
	"san/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
	userService *service.UserService
}

func NewPostHandler(postService *service.PostService, userService *service.UserService) *PostHandler {
	return &PostHandler{
		postService: postService,
		userService: userService,
	}
}

func (h *PostHandler) postToResponse(ctx context.Context, post *dbsqlc.Post) dto.PostResponse {
	user, err := h.userService.GetUserByID(ctx, post.UserID)
	var username, authorName string
	var avatarURL *string

	if err == nil && user != nil {
		username = user.Username
		authorName = user.Username // Fallback to username as name

		url, err := h.userService.GetAvatarURL(ctx, post.UserID)
		if err == nil && url != "" {
			avatarURL = &url
		}
	}

	var postImageURL *string
	url, err := h.postService.GetPostImageURL(ctx, post.ID)
	if err == nil && url != "" {
		postImageURL = &url
	}

	return dto.PostResponse{
		ID:              post.ID,
		UserID:          post.UserID,
		Title:           post.Title,
		Slug:            post.Slug,
		ImageURL:        postImageURL,
		Abstract:        post.Abstract,
		Body:            post.Body,
		Published:       post.Published,
		PublishDate:     post.PublishDate,
		Location:        post.Location,
		Lat:             post.Lat,
		Lon:             post.Lon,
		Locale:          post.Locale,
		Tags:            post.Tags,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
		AuthorUsername:  username,
		AuthorAvatarURL: avatarURL,
		AuthorName:      authorName,
		LikeCount:       0,     // TODO: Implement likes
		UserHasLiked:    false, // TODO: Implement likes
	}
}

// ListPostsByUserID godoc
// @Summary      List posts by user ID
// @Description  Get a list of posts for a specific user
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        id        path      string  true   "User ID"
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Success      200  {object}  response.Envelope{data=[]dto.PostResponse}
// @Router       /users/{id}/posts [get]
func (h *PostHandler) ListPostsByUserID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "me" {
		userID = c.MustGet("userID").(string)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, err := h.postService.ListPostsByUserID(c.Request.Context(), userID, int32(page), int32(pageSize))
	if err != nil {
		response.Error(c, err)
		return
	}

	var postResponses []dto.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, h.postToResponse(c.Request.Context(), post))
	}

	response.Success(c, http.StatusOK, postResponses)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	// Support both JSON and FormData?
	// The user asked for "attached image for post", implying file upload.
	// Typically file uploads use multipart/form-data.
	// So we should bind using ShouldBind which handles both if configured, or check Content-Type.
	// Given we expect a file, multipart is most likely.

	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, apperr.Unauthorized("User ID not found in context"))
		return
	}

	input := service.CreatePostInput{
		UserID:      userID,
		Title:       req.Title,
		Slug:        req.Slug,
		Abstract:    req.Abstract,
		Body:        req.Body,
		Published:   req.Published,
		PublishDate: req.PublishDate,
		Location:    req.Location,
		Lat:         req.Lat,
		Lon:         req.Lon,
		Locale:      req.Locale,
		Tags:        req.Tags,
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		input.ImageFile = file
		input.ImageSize = header.Size
		input.ImageContentType = header.Header.Get("Content-Type")
		input.ImageOriginalName = header.Filename
	}

	post, err := h.postService.CreatePost(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, h.postToResponse(c.Request.Context(), post))
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	post, err := h.postService.GetPostByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.postToResponse(c.Request.Context(), post))
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	posts, err := h.postService.ListPosts(c.Request.Context(), int32(page), int32(pageSize))
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = h.postToResponse(c.Request.Context(), post)
	}

	meta := response.MetaData{
		Page:     page,
		PageSize: pageSize,
		// TotalItems not implemented yet
	}

	response.SuccessWithMeta(c, http.StatusOK, responses, meta)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdatePostRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, apperr.Unauthorized("User ID not found in context"))
		return
	}

	input := service.UpdatePostInput{
		ID:          id,
		UserID:      userID,
		Title:       req.Title,
		Slug:        req.Slug,
		Abstract:    req.Abstract,
		Body:        req.Body,
		Published:   req.Published,
		PublishDate: req.PublishDate,
		Location:    req.Location,
		Lat:         req.Lat,
		Lon:         req.Lon,
		Locale:      req.Locale,
		Tags:        req.Tags,
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		input.ImageFile = file
		input.ImageSize = header.Size
		input.ImageContentType = header.Header.Get("Content-Type")
		input.ImageOriginalName = header.Filename
	}

	post, err := h.postService.UpdatePost(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.postToResponse(c.Request.Context(), post))
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		response.Error(c, apperr.Unauthorized("User ID not found in context"))
		return
	}

	err := h.postService.DeletePost(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
