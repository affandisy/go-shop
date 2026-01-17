package repository

import (
	"github.com/affandisy/goshop/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("User").Preload("OrderItems.Product.Category").Where("id = ?", id).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("User").Preload("OrderItems.Product.Category").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) GetByUserID(userID string, page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("OrderItems.Product").Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) GetAll(page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	if err := r.db.Model(&domain.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("User").Preload("OrderItems.Product").Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error
}
