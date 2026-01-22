package domain

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusExpired   PaymentStatus = "expired"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodGopay        PaymentMethod = "gopay"
	PaymentMethodShopeePay    PaymentMethod = "shopeepay"
	PaymentMethodQRIS         PaymentMethod = "qris"
	PaymentMethodAlfamart     PaymentMethod = "alfamart"
	PaymentMethodIndomaret    PaymentMethod = "indomaret"
)

type Payment struct {
	BaseModel
	OrderID           string        `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	Amount            float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentMethod     PaymentMethod `gorm:"type:varchar(50)" json:"payment_method"`
	Status            PaymentStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	MidtransOrderID   string        `gorm:"type:varchar(100);uniqueIndex" json:"midtrans_order_id"`
	MidtransSnapToken string        `gorm:"type:varchar(500)" json:"midtrans_snap_token"`
	MidtransSnapURL   string        `gorm:"type:varchar(500)" json:"midtrans_snap_url"`
	PaidAt            *time.Time    `json:"paid_at,omitempty"`
	ExpiredAt         *time.Time    `json:"expired_at,omitempty"`

	// Relasi
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

func (p *Payment) MarkAsPaid() {
	p.Status = PaymentStatusSuccess
	now := time.Now()
	p.PaidAt = &now
}

func (p *Payment) MarkAsFailed() {
	p.Status = PaymentStatusFailed
}

func (p *Payment) MarkAsExpired() {
	p.Status = PaymentStatusExpired
}
