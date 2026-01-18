package service

import (
	"context"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
	"github.com/affandisy/goshop/pkg/cache"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
	cacheService cache.CacheService
}

func NewCategoryService(categoryRepo repository.CategoryRepository, cacheService cache.CacheService) CategoryService {
	return &categoryService{categoryRepo: categoryRepo, cacheService: cacheService}
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

	ctx := context.Background()
	s.cacheService.Delete(ctx, cache.AllCategoriesKey())

	return category, nil
}

func (s *categoryService) GetByID(id string) (*domain.Category, error) {
	// return s.categoryRepo.GetByID(id)
	ctx := context.Background()
	cacheKey := cache.CategoryKey(id)

	var cachedCategory domain.Category
	err := s.cacheService.Get(ctx, cacheKey, &cachedCategory)
	if err == nil {
		return &cachedCategory, nil
	}

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.cacheService.Set(ctx, cacheKey, category, cache.CategoryTTL)

	return category, nil
}

func (s *categoryService) GetAll() ([]domain.Category, error) {
	// return s.categoryRepo.GetAll()
	ctx := context.Background()
	cacheKey := cache.AllCategoriesKey()

	var cachedCategories []domain.Category
	err := s.cacheService.Get(ctx, cacheKey, &cachedCategories)
	if err == nil {
		return cachedCategories, nil
	}

	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	s.cacheService.Set(ctx, cacheKey, categories, cache.CategoryTTL)

	return categories, nil
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

	ctx := context.Background()
	s.cacheService.Delete(ctx, cache.CategoryKey(id))
	s.cacheService.Delete(ctx, cache.AllCategoriesKey())
	s.cacheService.DeleteByPattern(ctx, cache.ProductsPrefix+"*")

	return category, nil
}

func (s *categoryService) Delete(id string) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	// return s.categoryRepo.Delete(id)

	err = s.categoryRepo.Delete(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	s.cacheService.Delete(ctx, cache.CategoryKey(id))
	s.cacheService.Delete(ctx, cache.AllCategoriesKey())
	s.cacheService.DeleteByPattern(ctx, cache.ProductsPrefix+"*")

	return nil
}
