package handler

import (
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/affandisy/goshop/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

var UpdateStockRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	product, err := h.productService.Create(req)
	if err != nil {
		response.InternalServerError(c, "Failed to create product", err)
		return
	}

	response.Created(c, "Product created successfully", product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalServerError(c, "Failed to get product", err)
		return
	}

	response.Success(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) List(c *gin.Context) {
	var query dto.ProductQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		query.Page = 1
		query.Limit = 10
	}

	products, total, err := h.productService.List(query)
	if err != nil {
		response.InternalServerError(c, "Failed to get products", err)
		return
	}

	paginationResp := utils.CreatePaginationResponse(query.Page, query.Limit, total, products)
	response.Success(c, "Products retrieved successfully", paginationResp)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	product, err := h.productService.Update(id, req)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product nof found")
			return
		}
		response.InternalServerError(c, "Failed to update product", err)
		return
	}

	response.Success(c, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.Delete(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalServerError(c, "Failed to delete product", err)
		return
	}

	response.Success(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("id")

	if err := c.ShouldBindJSON(&UpdateStockRequest); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	err := h.productService.UpdateStock(id, UpdateStockRequest.Quantity)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			response.BadRequest(c, "Insufficient stock", err)
			return
		}
		response.InternalServerError(c, "Failed to update stock", err)
		return
	}

	response.Success(c, "Stock updated successfully", nil)
}
