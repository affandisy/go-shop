package utils

import "gorm.io/gorm"

type PaginationParams struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

type PaginationResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if limit <= 0 || limit > 100 {
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func CreatePaginationResponse(page, limit int, totalRows int64, data interface{}) PaginationResponse {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	totalPages := int(totalRows) / limit
	if int(totalRows)%limit != 0 {
		totalPages++
	}

	return PaginationResponse{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
	}
}
