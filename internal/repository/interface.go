package repository

import (
	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	List(page, limit int) ([]domain.User, int64, error)
}

type CategoryRepository interface {
	Create(category *domain.Category) error
	GetByID(id string) (*domain.Category, error)
	GetAll() ([]domain.Category, error)
	Update(category *domain.Category) error
	Delete(id string) error
}

type ProductRepository interface {
	Create(product *domain.Product) error
	GetByID(id string) (*domain.Product, error)
	GetBySKU(sku string) (*domain.Product, error)
	List(query dto.ProductQuery) ([]domain.Product, int64, error)
	Update(product *domain.Product) error
	Delete(id string) error
	UpdateStock(id string, quantity int) error
}
