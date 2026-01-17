package repository

import (
	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Where("id = ?", id).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) GetBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) List(query dto.ProductQuery) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	db := r.db.Model(&domain.Product{}).Preload("Category")

	if query.Name != "" {
		db = db.Where("name ILIKE ?", "%"+query.Name+"%")
	}

	if query.CategoryID != "" {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	if query.MinPrice > 0 {
		db = db.Where("price >= ?", query.MinPrice)
	}

	if query.MaxPrice > 0 {
		db = db.Where("price >= 0", query.MinPrice)
	}

	if query.MaxPrice > 0 {
		db = db.Where("price <= ?", query.MaxPrice)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	offset := (query.Page - 1) * query.Limit
	err := db.Offset(offset).Limit(query.Limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id string) error {
	return r.db.Delete(&domain.Product{}, "id = ?", id).Error
}

func (r *productRepository) UpdateStock(id string, quantity int) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error
}
