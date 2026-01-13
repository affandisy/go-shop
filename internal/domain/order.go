package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	BaseModel
	OrderNumber string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID      string      `gorm:"type:uuid;not null" json:"user_id"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      OrderStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Notes       string      `gorm:"type:text" json:"notes"`
	PaidAt      *time.Time  `json:"paid_at,omitempty"`

	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) CalculateTotalAmount() {
	total := 0.0
	for _, item := range o.OrderItems {
		total += item.Price * float64(item.Quantity)
	}
	o.TotalAmount = total
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusPaid
}

func (o *Order) MarkAsPaid() {
	o.Status = OrderStatusPaid
	now := time.Now()
	o.PaidAt = &now
}
