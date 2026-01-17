package service

import (
	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(req dto.CategoryRequest) (*domain.Category, error) {
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetByID(id string) (*domain.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetAll() ([]domain.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) Update(id string, req dto.CategoryRequest) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id string) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(id)
}
