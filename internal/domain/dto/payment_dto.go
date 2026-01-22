package dto

import (
	"time"

	"github.com/affandisy/goshop/internal/domain"
)

type CreatePaymentRequest struct {
	OrderID       string               `json:"order_id" binding:"required"`
	PaymentMethod domain.PaymentMethod `json:"payment_method"`
}

type PaymentNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionID     string `json:"transaction_id"`
	FraudStatus       string `json:"fraud_status"`
	SignatureKey      string `json:"signature_key"`
}

type MidtransResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type PaymentResponse struct {
	ID                string               `json:"id"`
	OrderID           string               `json:"order_id"`
	OrderNumber       string               `json:"order_number"`
	Amount            float64              `json:"amount"`
	PaymentMethod     domain.PaymentMethod `json:"payment_method"`
	Status            domain.PaymentStatus `json:"status"`
	MidtransSnapToken string               `json:"snap_token,omitempty"`
	MidtransSnapURL   string               `json:"snap_url,omitempty"`
	ExpiredAt         *time.Time           `json:"expired_at,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
}

func PaymentMapToResponse(p *domain.Payment) PaymentResponse {
	orderNumber := ""
	if p.Order != nil {
		orderNumber = p.Order.OrderNumber
	}
	return PaymentResponse{
		ID:                p.ID,
		OrderID:           p.OrderID,
		OrderNumber:       orderNumber,
		Amount:            p.Amount,
		PaymentMethod:     p.PaymentMethod,
		Status:            p.Status,
		MidtransSnapToken: p.MidtransSnapToken,
		MidtransSnapURL:   p.MidtransSnapURL,
		ExpiredAt:         p.ExpiredAt,
		CreatedAt:         p.CreatedAt,
	}
}
