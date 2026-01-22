package handler

import (
	"errors"
	"log"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/middleware"
	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/affandisy/goshop/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	payment, err := h.paymentService.CreatePayment(userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.NotFound(c, "Order not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			response.Forbidden(c, "You don't have access to this order")
			return
		}
		if errors.Is(err, domain.ErrOrderAlreadyPaid) {
			response.BadRequest(c, "Order already paid", err)
			return
		}
		if errors.Is(err, domain.ErrPaymentAlreadyExists) {
			response.BadRequest(c, "Payment already exists for this order", err)
			return
		}
		response.InternalServerError(c, "Failed to create payment", err)
		return
	}

	response.Created(c, "Payment created successfully. Please complete payment via Snap URL", dto.PaymentMapToResponse(payment))
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	id := c.Param("id")

	payment, err := h.paymentService.GetPaymentByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			response.NotFound(c, "Payment not found")
			return
		}
		response.InternalServerError(c, "Failed to get payment", err)
		return
	}

	response.Success(c, "Payment retrieved successfully", dto.PaymentMapToResponse(payment))
}

func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")

	payment, err := h.paymentService.GetPaymentByOrderID(orderID)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			response.NotFound(c, "Payment not found")
			return
		}
		response.InternalServerError(c, "Failed to get payment", err)
		return
	}

	response.Success(c, "Payment retrieved sucessfully", dto.PaymentMapToResponse(payment))
}

func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	var notification dto.PaymentNotification

	if err := c.ShouldBindJSON(&notification); err != nil {
		log.Printf("Invalid notification body: %v", err)
		response.BadRequest(c, "Invalid notification body", err)
		return
	}

	log.Printf("Received payment notification for order: %s, status: %s",
		notification.OrderID, notification.TransactionStatus)

	err := h.paymentService.HandleNotification(notification)
	if err != nil {
		log.Printf("Failed to handle notification: %v", err)
		response.InternalServerError(c, "Failed to process notification", err)
		return
	}

	response.Success(c, "Notification processed successfully", nil)
}

func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	var params utils.PaginationParams

	if err := c.ShouldBindQuery(&params); err != nil {
		params.Page = 1
		params.Limit = 10
	}

	payments, total, err := h.paymentService.GetAllPayments(params.Page, params.Limit)
	if err != nil {
		response.InternalServerError(c, "Failed to get payments", err)
		return
	}

	paymentResponses := make([]dto.PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = dto.PaymentMapToResponse(&payment)
	}

	paginationResp := utils.CreatePaginationResponse(params.Page, params.Limit, total, paymentResponses)
	response.Success(c, "Payments retrieved successfully", paginationResp)
}
