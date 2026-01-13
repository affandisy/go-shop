package dto

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	SKU         string  `json:"sku" binding:"required"`
	CategoryID  string  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

type ProductQuery struct {
	Name       string  `form:"name"`
	CategoryID string  `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	IsActive   *bool   `form:"is_active"`
	Page       int     `form:"page,default=1"`
	Limit      int     `form:"limit,default=10"`
}
