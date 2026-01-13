package domain

type OrderItem struct {
	BaseModel
	OrderID   string  `gorm:"type:uuid;not null" json:"order_id"`
	ProductID string  `gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"` // harga saat order dibuat

	// Relasi
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
