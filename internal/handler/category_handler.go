package handler

import (
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	category, err := h.categoryService.Create(req)
	if err != nil {
		response.InternalServerError(c, "Failed to create category", err)
		return
	}

	response.Created(c, "Category created successfully", category)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalServerError(c, "Failed to get category", err)
		return
	}

	response.Success(c, "Category retrieved successfully", category)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		response.InternalServerError(c, "Failed to get categories", err)
		return
	}

	response.Success(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	category, err := h.categoryService.Update(id, req)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalServerError(c, "Failed to update category", err)
		return
	}

	response.Success(c, "Category updated successfully", category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryService.Delete(id)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalServerError(c, "Failed to delete category", err)
		return
	}

	response.Success(c, "Category deleted successfully", nil)
}
