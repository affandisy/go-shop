package handler

import (
	"context"

	"github.com/affandisy/goshop/pkg/cache"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/gin-gonic/gin"
)

type CacheHandler struct {
	cacheService cache.CacheService
}

func NewCacheHandler(cacheService cache.CacheService) *CacheHandler {
	return &CacheHandler{cacheService: cacheService}
}

func (h *CacheHandler) ClearProductsCache(c *gin.Context) {
	ctx := context.Background()

	err := h.cacheService.DeleteByPattern(ctx, cache.ProductsPrefix+"*")
	if err != nil {
		response.InternalServerError(c, "Failed to clear cache", err)
		return
	}

	err = h.cacheService.DeleteByPattern(ctx, cache.ProductPrefix+"*")
	if err != nil {
		response.InternalServerError(c, "Failed to clear cache", err)
		return
	}

	response.Success(c, "Products cache cleared successfully", nil)
}

func (h *CacheHandler) ClearCategoriesCache(c *gin.Context) {
	ctx := context.Background()

	err := h.cacheService.DeleteByPattern(ctx, cache.CategoryPrefix+"*")
	if err != nil {
		response.InternalServerError(c, "Failed to clear cache", err)
		return
	}

	response.Success(c, "Categories cache cleared successfully", nil)
}

func (h *CacheHandler) ClearAllCache(c *gin.Context) {
	ctx := context.Background()

	err := h.cacheService.FlushAll(ctx)
	if err != nil {
		response.InternalServerError(c, "Failed to clear cache", err)
		return
	}

	response.Success(c, "All cache cleared successfully", nil)
}

func (h *CacheHandler) GetCacheStats(c *gin.Context) {
	ctx := context.Background()

	stats := gin.H{
		"products_cached":   h.cacheService.Exists(ctx, cache.ProductsPrefix+"list:page:1:limit:10"),
		"categories_cached": h.cacheService.Exists(ctx, cache.AllCategoriesKey()),
		"cache_enabled":     true,
	}

	response.Success(c, "Cache statistics retrieved", stats)
}
