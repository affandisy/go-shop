package cache

import (
	"fmt"
	"time"
)

const (
	ProductPrefix  = "product:"
	ProductsPrefix = "products:"
	CategoryPrefix = "category:"
	UserPrefix     = "user:"
	OrderPrefix    = "order:"
)

const (
	ProductTTL  = 15 * time.Minute // Product cache 15 menit
	CategoryTTL = 30 * time.Minute // Category cache 30 menit
	ProductsTTL = 5 * time.Minute  // Products list cache 5 menit
	UserTTL     = 10 * time.Minute // User cache 10 menit
	OrderTTL    = 5 * time.Minute  // Order cache 5 menit
)

func ProductKey(id string) string {
	return fmt.Sprintf("%s%s", ProductPrefix, id)
}

func ProductsKey(page, limit int, filters map[string]interface{}) string {
	key := fmt.Sprintf("%slist:page:%d:limit:%d", ProductPrefix, page, limit)

	if name, ok := filters["name"].(string); ok && name != "" {
		key += fmt.Sprintf(":name:%s", name)
	}
	if categoryID, ok := filters["category_id"].(string); ok && categoryID != "" {
		key += fmt.Sprintf(":cat:%s", categoryID)
	}
	if minPrice, ok := filters["min_price"].(float64); ok && minPrice > 0 {
		key += fmt.Sprintf(":minp:%.0f", minPrice)
	}
	if maxPrice, ok := filters["max_price"].(float64); ok && maxPrice > 0 {
		key += fmt.Sprintf(":maxp:%.0f", maxPrice)
	}

	return key
}

func CategoryKey(id string) string {
	return fmt.Sprintf("%s%s", CategoryPrefix, id)
}

func AllCategoriesKey() string {
	return fmt.Sprintf("%sall", CategoryPrefix)
}

func UserKey(id string) string {
	return fmt.Sprintf("%s%s", UserPrefix, id)
}

func OrderKey(id string) string {
	return fmt.Sprintf("%s%s", OrderPrefix, id)
}

func UserOrdersKey(userID string, page, limit int) string {
	return fmt.Sprintf("%suser:%s:page:%d:limit:%d", OrderPrefix, userID, page, limit)
}
