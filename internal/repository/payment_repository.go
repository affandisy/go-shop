package repository

import (
	"github.com/affandisy/goshop/internal/domain"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByID(id string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Preload("Order.User").Where("id = ?", id).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) GetByOrderID(orderID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Preload("Order.User").Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) GetByMidtransOrderID(midtransOrderID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Preload("Order.User").Where("midtrans_order_id = ?", midtransOrderID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) Update(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) List(page, limit int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	if err := r.db.Model(&domain.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("Order.User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
