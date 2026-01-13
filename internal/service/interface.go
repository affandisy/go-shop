package service

import (
	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
)

type UserService interface {
	Register(req dto.UserRegisterRequest) (*domain.User, error)
	Login(req dto.UserLoginRequest) (*domain.User, string, error)
	GetProfile(userID string) (*domain.User, error)
	UpdateProfile(userID string, req dto.UserRegisterRequest) (*domain.User, error)
	GetUsers(page, limit int) ([]domain.User, int64, error)
}
