package handler

import (
	"errors"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/middleware"
	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/affandisy/goshop/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	order, err := h.orderService.CreateOrder(userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyCart) {
			response.BadRequest(c, "Cart is empty", err)
			return
		}
		if errors.Is(err, domain.ErrProductNotFound) {
			response.BadRequest(c, "Product not found", err)
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			response.BadRequest(c, "Insufficient stock", err)
			return
		}
		response.InternalServerError(c, "Failed to create order", err)
		return
	}

	response.Created(c, "Order created successfully", order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	role, _ := middleware.GetUserRole(c)
	isAdmin := role == "admin"

	orderID := c.Param("id")

	order, err := h.orderService.GetOrderByID(orderID, userID, isAdmin)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.NotFound(c, "Order not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Forbidden(c, "You dont have access to this order")
			return
		}
		response.InternalServerError(c, "Failed to get order", err)
		return
	}

	response.Success(c, "Order retrieved successfully", order)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params.Page = 1
		params.Limit = 10
	}

	orders, total, err := h.orderService.GetMyOrders(userID, params.Page, params.Limit)
	if err != nil {
		response.InternalServerError(c, "Failed to get orders", err)
		return
	}

	paginationResp := utils.CreatePaginationResponse(params.Page, params.Limit, total, orders)
	response.Success(c, "Orders retrieved successfully", paginationResp)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params.Page = 1
		params.Limit = 10
	}

	orders, total, err := h.orderService.GetAllOrders(params.Page, params.Limit)
	if err != nil {
		response.InternalServerError(c, "Failed to get orders", err)
		return
	}

	paginationResp := utils.CreatePaginationResponse(params.Page, params.Limit, total, orders)
	response.Success(c, "Orders retrieved successfully", paginationResp)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	order, err := h.orderService.UpdateOrderStatus(orderID, req)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.NotFound(c, "Order not found")
			return
		}
		response.InternalServerError(c, "Failed to update order status", err)
		return
	}

	response.Success(c, "Order status updated successfully", order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authentication")
		return
	}

	orderID := c.Param("id")

	err = h.orderService.CancelOrder(orderID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.NotFound(c, "Order not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Forbidden(c, "You dont have access to this order")
			return
		}
		if errors.Is(err, domain.ErrCannotCancelOrder) {
			response.BadRequest(c, "Order cannot be cancelled", err)
			return
		}
		response.InternalServerError(c, "Failed to cancel order", err)
		return
	}

	response.Success(c, "Order cancelled successfully", nil)
}
