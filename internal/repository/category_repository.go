package repository

import (
	"github.com/affandisy/goshop/internal/domain"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Where("is_active = ?", true).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id string) error {
	return r.db.Delete(&domain.Category{}, "id = ?", id).Error
}
