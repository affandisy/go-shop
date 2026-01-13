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
	}

	log.Println("Seeding completed successfully")

	return nil
}
