package handler

import (
	"context"
	"net/http"
	"strconv"

	"san/api/dto"
	dbsqlc "san/internal/db/sqlc"
	"san/internal/service"
	"san/pkg/apperr"
	"san/pkg/response"

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
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, h.userToResponse(c.Request.Context(), user))
}

// GetUserByID godoc
// @Summary      Get a user by ID
// @Description  Get details of a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Envelope{data=dto.UserResponse}
// @Failure      404  {object}  response.Envelope{error=response.ErrorData}
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.userToResponse(ctx, user))
}

// ListUsers godoc
// @Summary      List all users
// @Description  Get a list of all users with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number (default 1)"
// @Param        page_size query     int     false  "Page size (default 10)"
// @Success      200  {object}  response.Envelope{data=[]dto.UserResponse,meta=response.MetaData}
// @Failure      500  {object}  response.Envelope{error=response.ErrorData}
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.service.ListUsers(ctx, int32(page), int32(pageSize))
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, h.userToResponse(ctx, u))
	}

	meta := response.MetaData{
		Page:     page,
		PageSize: pageSize,
		// TotalItems and TotalPages should be fetched from service if available
	}

	response.SuccessWithMeta(c, http.StatusOK, responses, meta)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update user details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      string  true   "User ID"
// @Param        username body      string  false  "Username"
// @Param        email    body      string  false  "Email"
// @Param        bio      body      string  false  "Bio"
// @Success      200  {object}  response.Envelope{data=dto.UserResponse}
// @Failure      400  {object}  response.Envelope{error=response.ErrorData}
// @Failure      404  {object}  response.Envelope{error=response.ErrorData}
// @Failure      500  {object}  response.Envelope{error=response.ErrorData}
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Bio      *string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	ctx := c.Request.Context()
	input := service.UpdateUserInput{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Bio:      req.Bio,
	}

	user, err := h.service.UpdateUser(ctx, input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.userToResponse(ctx, user))
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope{error=response.ErrorData}
// @Failure      500  {object}  response.Envelope{error=response.ErrorData}
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	if err := h.service.DeleteUser(ctx, id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// UploadAvatar godoc
// @Summary      Upload user avatar
// @Description  Upload an avatar image for a user
// @Tags         users
// @Accept       mpfd
// @Produce      json
// @Param        id      path      string  true  "User ID"
// @Param        avatar  formData  file    true  "Avatar Image"
// @Success      200  {object}  response.Envelope{data=dto.UserResponse}
// @Failure      400  {object}  response.Envelope{error=response.ErrorData}
// @Failure      500  {object}  response.Envelope{error=response.ErrorData}
// @Router       /users/{id}/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.Error(c, apperr.BadRequest("user id is required"))
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.Error(c, apperr.BadRequest("failed to get file: "+err.Error()))
		return
	}
	defer file.Close()

	ctx := c.Request.Context()
	user, err := h.service.UploadUserAvatar(ctx, userID, file, header.Size, header.Header.Get("Content-Type"), header.Filename)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.userToResponse(ctx, user))
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
