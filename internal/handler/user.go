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

	response.Success(c, http.StatusOK, h.userToResponse(c.Request.Context(), user))
}

// ListUsers godoc
// @Summary      List users
// @Description  Get a list of users with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Success      200  {object}  response.Envelope{data=[]dto.UserResponse}
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	ctx := c.Request.Context()

	users, err := h.service.ListUsers(ctx, int32(page), int32(pageSize))
	if err != nil {
		response.Error(c, err)
		return
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, h.userToResponse(c.Request.Context(), user))
	}

	response.Success(c, http.StatusOK, userResponses)
}

// VerifyEmail godoc
// @Summary      Verify email with OTP
// @Description  Verify user email address using OTP
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.VerifyEmailRequest true "Verification Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/verify [post]
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), req.Email, req.OTP); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// Login godoc
// @Summary      User login
// @Description  Login with username/email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login Request"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	result, err := h.service.Login(c.Request.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         h.userToResponse(c.Request.Context(), result.User),
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success      200  {object}  dto.TokenResponse
// @Failure      400  {object}  map[string]string
// @Router       /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user profile
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path     string                 true  "User ID"
// @Param        request body     dto.UpdateUserRequest  true  "Update Request"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperr.BadRequest(err.Error()))
		return
	}

	input := service.UpdateUserInput{
		ID:       id,
		Username: req.Username,
		Bio:      req.Bio,
		Email:    req.Email,
	}

	user, err := h.service.UpdateUser(c.Request.Context(), input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.userToResponse(c.Request.Context(), user))
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UploadAvatar godoc
// @Summary      Upload user avatar
// @Description  Upload a new avatar for the user
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id      path      string  true  "User ID"
// @Param        avatar  formData  file    true  "Avatar Image"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Router       /users/{id}/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	id := c.Param("id")
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.Error(c, apperr.BadRequest("Avatar file is required"))
		return
	}
	defer file.Close()

	user, err := h.service.UploadUserAvatar(c.Request.Context(), id, file, header.Size, header.Header.Get("Content-Type"), header.Filename)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, h.userToResponse(c.Request.Context(), user))
}

func (h *UserHandler) userToResponse(ctx context.Context, user *dbsqlc.User) dto.UserResponse {
	avatarURL, _ := h.service.GetAvatarURL(ctx, user.ID)

	var image *string
	if avatarURL != "" {
		image = &avatarURL
	}

	return dto.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Bio:        user.Bio,
		Image:      image,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
