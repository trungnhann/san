package handler

import (
	"context"
	"net/http"

	"san/api/dto"
	dbsqlc "san/internal/db/sqlc"
	"san/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserUseCase
}

func NewUserHandler(service UserUseCase) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with optional avatar
// @Tags         users
// @Accept       json,mpfd
// @Produce      json
// @Param        username  formData  string  true   "Username"
// @Param        email     formData  string  true   "Email"
// @Param        password  formData  string  true   "Password"
// @Param        bio       formData  string  false  "Bio"
// @Param        avatar    formData  file    false  "Avatar Image"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Prepare service input
	input := service.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Bio:      req.Bio,
	}

	file, header, err := c.Request.FormFile("avatar")
	if err == nil {
		defer file.Close()
		input.AvatarFile = file
		input.AvatarSize = header.Size
		input.AvatarContentType = header.Header.Get("Content-Type")
		input.AvatarOriginalName = header.Filename
	}

	ctx := c.Request.Context()

	user, err := h.service.CreateUser(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, h.userToResponse(c.Request.Context(), user))
}

// GetUserByID godoc
// @Summary      Get a user by ID
// @Description  Get details of a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      404  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, h.userToResponse(ctx, user))
}

// ListUsers godoc
// @Summary      List all users
// @Description  Get a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.UserResponse
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.service.ListUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, h.userToResponse(ctx, u))
	}

	c.JSON(http.StatusOK, responses)
}

// UploadAvatar godoc
// @Summary      Upload user avatar
// @Description  Upload an avatar image for a user
// @Tags         users
// @Accept       mpfd
// @Produce      json
// @Param        id      path      string  true  "User ID"
// @Param        avatar  formData  file    true  "Avatar Image"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id}/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file: " + err.Error()})
		return
	}
	defer file.Close()

	ctx := c.Request.Context()
	user, err := h.service.UploadUserAvatar(ctx, userID, file, header.Size, header.Header.Get("Content-Type"), header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload avatar"})
		return
	}

	c.JSON(http.StatusOK, h.userToResponse(ctx, user))
}

func (h *UserHandler) userToResponse(ctx context.Context, u *dbsqlc.User) dto.UserResponse {
	avatarURL := ""
	if url, err := h.service.GetAvatarURL(ctx, u.ID); err == nil {
		avatarURL = url
	}

	imagePtr := &avatarURL
	if avatarURL == "" {
		imagePtr = nil
	}

	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Bio:       u.Bio,
		Image:     imagePtr,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
