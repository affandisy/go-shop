package service

import (
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
)

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{productRepo: productRepo, categoryRepo: categoryRepo}
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
	return s.productRepo.GetByID(id)
}

func (s *productService) List(query dto.ProductQuery) ([]domain.Product, int64, error) {
	return s.productRepo.List(query)
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

	product, _ = s.productRepo.GetByID(product.ID)

	return product, nil
}

func (s *productService) Delete(id string) error {
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}

func (s *productService) UpdateStock(id string, quantity int) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	if product.Stock+quantity < 0 {
		return domain.ErrInsufficientStock
	}

	return s.productRepo.UpdateStock(id, quantity)
}
