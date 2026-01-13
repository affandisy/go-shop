package domain

type Product struct {
	BaseModel
	Name        string  `gorm:"type:varchar(200);not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int     `gorm:"not null;default:0" json:"stock"`
	SKU         string  `gorm:"type:varchar(100);uniqueIndex" json:"sku"` // Stock Keeping Unit
	CategoryID  string  `gorm:"type:uuid" json:"category_id"`
	ImageURL    string  `gorm:"type:varchar(500)" json:"image_url"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) IsAvailable() bool {
	return p.IsActive && p.Stock > 0
}

func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return ErrInsufficientStock
	}
	p.Stock -= quantity
	return nil
}
