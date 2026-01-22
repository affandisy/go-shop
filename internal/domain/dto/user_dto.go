package dto

import "github.com/affandisy/goshop/internal/domain"

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

func UserMapToResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Name:     u.Name,
		Phone:    u.Phone,
		Role:     u.Role,
		IsActive: u.IsActive,
	}
}
