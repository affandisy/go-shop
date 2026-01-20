package service

import (
	"time"

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

type CategoryService interface {
	Create(req dto.CategoryRequest) (*domain.Category, error)
	GetByID(id string) (*domain.Category, error)
	GetAll() ([]domain.Category, error)
	Update(id string, req dto.CategoryRequest) (*domain.Category, error)
	Delete(id string) error
}

type ProductService interface {
	Create(req dto.ProductRequest) (*domain.Product, error)
	GetByID(id string) (*domain.Product, error)
	List(query dto.ProductQuery) ([]domain.Product, int64, error)
	Update(id string, req dto.ProductRequest) (*domain.Product, error)
	Delete(id string) error
	UpdateStock(id string, quantity int) error
}

type OrderService interface {
	CreateOrder(userID string, req dto.CreateOrderRequest) (*domain.Order, error)
	GetOrderByID(orderID string, userID string, isAdmin bool) (*domain.Order, error)
	GetMyOrders(userID string, page, limit int) ([]domain.Order, int64, error)
	GetAllOrders(page, limit int) ([]domain.Order, int64, error)
	UpdateOrderStatus(orderID string, req dto.UpdateOrderStatusRequest) (*domain.Order, error)
	CancelOrder(orderID string, userID string) error
}

type ReportService interface {
	GenerateUsersReport(format string) ([]byte, string, error)
	GenerateProductsReport(format string) ([]byte, string, error)
	GenerateOrdersReport(startDate, endDate time.Time, format string) ([]byte, string, error)
}
