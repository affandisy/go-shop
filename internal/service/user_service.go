package service

import (
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
	"github.com/affandisy/goshop/pkg/utils"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(req dto.UserRegisterRequest) (*domain.User, error) {
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &domain.User{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     "customer",
		IsActive: true,
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req dto.UserLoginRequest) (*domain.User, string, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, "", domain.ErrInvalidCredentials
		}

		return nil, "", err
	}

	if !user.IsActive {
		return nil, "", domain.ErrUnauthorized
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, "", domain.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *userService) GetProfile(userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(userID string, req dto.UserRegisterRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Phone = req.Phone

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUsers(page, limit int) ([]domain.User, int64, error) {
	users, total, err := s.userRepo.List(page, limit)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
