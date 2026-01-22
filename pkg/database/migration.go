package database

import (
	"log"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/pkg/utils"
)

func AutoMigrate() error {
	log.Println("Running auto migration...")

	err := DB.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)

	if err != nil {
		return err
	}

	log.Println("Auto migration completed successfully")
	return nil
}

func SeedData() error {
	log.Println("Seeding initial data...")

	var count int64
	DB.Model(&domain.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		admin := &domain.User{
			Email:    "admin@goshop.com",
			Name:     "Admin",
			Password: "admin123",
			Phone:    "081234567890",
			Role:     "admin",
			IsActive: true,
		}

		hashedPassword, err := utils.HashPassword(admin.Password)
		if err != nil {
			return err
		}

		admin.Password = hashedPassword

		if err := DB.Create(admin).Error; err != nil {
			return err
		}

		log.Println("Default admin created: %w, %w", admin.Email, admin.Password)
	}

	DB.Model(&domain.Category{}).Count(&count)
	if count == 0 {
		categories := []domain.Category{
			{Name: "Electronics", Description: "Electronic gadgets and devices"},
			{Name: "Books", Description: "Various kinds of books"},
			{Name: "Clothing", Description: "Apparel and accessories"},
			{Name: "Home & Living", Description: "Home decoration and furniture"},
		}

		if err := DB.Create(&categories).Error; err != nil {
			return err
		}

		log.Println("Default categories created")

		var electronicsCategory domain.Category
		DB.Where("name = ?", "Electronics").First(&electronicsCategory)

		var smartphonesCategory domain.Category
		DB.Where("name = ?", "Smartphones").First(&smartphonesCategory)

		products := []domain.Product{
			{
				Name:        "MacBook Pro 16 M3 Max",
				Description: "Professional laptop with M3 Max chip, 36GB RAM, 1TB SSD",
				Price:       49999000,
				Stock:       10,
				SKU:         "MBP-16-M3MAX-1TB",
				CategoryID:  electronicsCategory.ID,
				ImageURL:    "https://example.com/macbook-pro.jpg",
				IsActive:    true,
			},
			{
				Name:        "iPhone 15 Pro Max",
				Description: "Latest iPhone with A17 Pro chip, titanium design, 256GB",
				Price:       19999000,
				Stock:       50,
				SKU:         "IPH-15-PM-256-BLU",
				CategoryID:  smartphonesCategory.ID,
				ImageURL:    "https://example.com/iphone-15-pro.jpg",
				IsActive:    true,
			},
			{
				Name:        "Samsung Galaxy S24 Ultra",
				Description: "Premium Android flagship with S Pen, 512GB storage",
				Price:       18999000,
				Stock:       30,
				SKU:         "SAM-S24-ULTRA-512-BLK",
				CategoryID:  smartphonesCategory.ID,
				ImageURL:    "https://example.com/galaxy-s24.jpg",
				IsActive:    true,
			},
		}

		if err := DB.Create(&products).Error; err != nil {
			return err
		}

		log.Println("Default products created")
	}

	log.Println("Seeding completed successfully")

	return nil
}
