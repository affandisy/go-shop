package dto

import "github.com/affandisy/goshop/internal/domain"

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
	Notes string             `json:"notes"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status" binding:"required,oneof=pending paid processing shipped delivered cancelled"`
}
