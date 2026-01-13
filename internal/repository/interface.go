package repository

import "github.com/affandisy/goshop/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	List(page, limit int) ([]domain.User, int64, error)
}
