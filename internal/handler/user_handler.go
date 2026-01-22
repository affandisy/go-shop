package handler

import (
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/middleware"
	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/affandisy/goshop/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			response.BadRequest(c, "Email already registered", err)
			return
		}

		response.InternalServerError(c, "Failed to register user", err)
		return
	}

	response.Created(c, "User registered successfully", dto.UserMapToResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	user, token, err := h.userService.Login(req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			response.Unauthorized(c, "Unauthorized access")
			return
		}
		response.InternalServerError(c, "Failed to login", err)
		return
	}

	response.Success(c, "Login Successful", gin.H{
		"user":  dto.UserMapToResponse(user),
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Unauthorized access")
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}

		response.InternalServerError(c, "Failed to get user profile", err)
		return
	}

	response.Success(c, "Profile retrieved successfully", dto.UserMapToResponse(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "Unauthorized access")
		return
	}

	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	user, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to update profile", err)
		return
	}

	response.Success(c, "Profile updated successfully", dto.UserMapToResponse(user))
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params.Page = 1
		params.Limit = 10
	}

	users, total, err := h.userService.GetUsers(params.Page, params.Limit)
	if err != nil {
		response.InternalServerError(c, "Failed to get users", err)
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserMapToResponse(&user)
	}

	paginationResp := utils.CreatePaginationResponse(params.Page, params.Limit, total, userResponses)

	response.Success(c, "Users retrieved successfully", paginationResp)
}
