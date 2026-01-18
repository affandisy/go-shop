package service

import (
	"context"
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
	"github.com/affandisy/goshop/pkg/cache"
)

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	cacheService cache.CacheService
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, cacheService cache.CacheService) ProductService {
	return &productService{productRepo: productRepo, categoryRepo: categoryRepo, cacheService: cacheService}
}

func (s *productService) Create(req dto.ProductRequest) (*domain.Product, error) {
	existingProduct, err := s.productRepo.GetBySKU(req.SKU)
	if err != nil && !errors.Is(err, domain.ErrProductNotFound) {
		return nil, err
	}

	if existingProduct != nil {
		return nil, errors.New("sku already exists")
	}

	if req.CategoryID != "" {
		_, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
		CategoryID:  req.CategoryID,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	product, _ = s.productRepo.GetByID(product.ID)

	return product, nil
}

func (s *productService) GetByID(id string) (*domain.Product, error) {
	// return s.productRepo.GetByID(id)
	ctx := context.Background()
	cacheKey := cache.ProductKey(id)

	var cachedProduct domain.Product
	err := s.cacheService.Get(ctx, cacheKey, &cachedProduct)
	if err == nil {
		return &cachedProduct, nil
	}

	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.cacheService.Set(ctx, cacheKey, product, cache.ProductTTL)

	return product, nil
}

func (s *productService) List(query dto.ProductQuery) ([]domain.Product, int64, error) {
	// return s.productRepo.List(query)
	ctx := context.Background()

	filters := map[string]interface{}{
		"name":        query.Name,
		"category_id": query.CategoryID,
		"min_price":   query.MinPrice,
		"max_price":   query.MaxPrice,
	}
	cacheKey := cache.ProductsKey(query.Page, query.Limit, filters)

	var cachedResult struct {
		Products   []domain.Product `json:"products"`
		TotalCount int64            `json:"total_count"`
	}

	err := s.cacheService.Get(ctx, cacheKey, &cachedResult)
	if err == nil {
		return cachedResult.Products, cachedResult.TotalCount, nil
	}

	products, total, err := s.productRepo.List(query)
	if err != nil {
		return nil, 0, err
	}

	cachedResult.Products = products
	cachedResult.TotalCount = total
	s.cacheService.Set(ctx, cacheKey, cachedResult, cache.ProductTTL)

	return products, total, nil
}

func (s *productService) Update(id string, req dto.ProductRequest) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.SKU != product.SKU {
		existingProduct, err := s.productRepo.GetBySKU(req.SKU)
		if err != nil && !errors.Is(err, domain.ErrProductNotFound) {
			return nil, err
		}
		if existingProduct != nil {
			return nil, errors.New("sku already exists")
		}
	}

	if req.CategoryID != "" && req.CategoryID != product.CategoryID {
		_, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return nil, err
		}
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.SKU = req.SKU
	product.CategoryID = req.CategoryID
	product.ImageURL = req.ImageURL

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	ctx := context.Background()
	s.cacheService.Delete(ctx, cache.ProductKey(id))
	s.cacheService.DeleteByPattern(ctx, cache.ProductsPrefix+"*")

	product, _ = s.productRepo.GetByID(product.ID)

	s.cacheService.Set(ctx, cache.ProductKey(id), product, cache.ProductTTL)

	return product, nil
}

func (s *productService) Delete(id string) error {
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	err = s.productRepo.Delete(id)
	if err != nil {
		return err
	}

	ctx := context.Background()
	s.cacheService.Delete(ctx, cache.ProductKey(id))
	s.cacheService.DeleteByPattern(ctx, cache.ProductsPrefix+"*")

	return nil
}

func (s *productService) UpdateStock(id string, quantity int) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	if product.Stock+quantity < 0 {
		return domain.ErrInsufficientStock
	}

	err = s.productRepo.UpdateStock(id, quantity)
	if err != nil {
		return err
	}

	ctx := context.Background()
	s.cacheService.Delete(ctx, cache.ProductKey(id))
	s.cacheService.DeleteByPattern(ctx, cache.ProductsPrefix+"*")

	return nil
}
