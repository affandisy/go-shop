package service

import (
	"fmt"
	"time"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
	"github.com/affandisy/goshop/pkg/utils"
)

type reportService struct {
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
}

func NewReportService(userRepo repository.UserRepository, productRepo repository.ProductRepository, orderRepo repository.OrderRepository) ReportService {
	return &reportService{
		userRepo:    userRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

func (s *reportService) GenerateUsersReport(format string) ([]byte, string, error) {
	users, _, err := s.userRepo.List(1, 1000)
	if err != nil {
		return nil, "", err
	}

	if format == "pdf" {
		return s.generateUsersPDF(users)
	}

	return s.generateUsersExcel(users)
}

func (s *reportService) generateUsersPDF(users []domain.User) ([]byte, string, error) {
	pdf := utils.NewPDFGenerator()
	pdf.SetTitle("Users Report")

	headers := []string{"No", "Name", "Email", "Phone", "Role", "Status"}
	widths := []float64{10, 50, 60, 35, 25, 20}
	pdf.AddTableHeader(headers, widths)

	for i, user := range users {
		status := "Active"
		if !user.IsActive {
			status = "Inactive"
		}

		values := []string{
			fmt.Sprintf("%d", i+1),
			user.Name,
			user.Email,
			user.Phone,
			user.Role,
			status,
		}
		pdf.AddTableRow(values, widths)
	}

	summary := map[string]string{
		"Total Users": fmt.Sprintf("%d", len(users)),
		"Report Date": time.Now().Format("2006-01-02"),
	}
	pdf.AddSummary("Summary", summary)

	data, err := pdf.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("users_report_%s.pdf", time.Now().Format("20060102_150405"))
	return data, filename, nil
}

func (s *reportService) generateUsersExcel(users []domain.User) ([]byte, string, error) {
	excel := utils.NewExcelGenerator("Users Report")
	excel.SetTitle("Users Report - " + time.Now().Format("2006-01-02"))

	headers := []string{"No", "Name", "Email", "Phone", "Role", "Status", "Created At"}
	excel.AddTableHeader(headers)

	for i, user := range users {
		status := "Active"
		if !user.IsActive {
			status = "Inactive"
		}

		values := []interface{}{
			i + 1,
			user.Name,
			user.Email,
			user.Phone,
			user.Role,
			status,
			user.CreatedAt.Format("2006-01-02"),
		}
		excel.AddTableRow(values)
	}

	summary := map[string]interface{}{
		"Total Users": len(users),
		"Report Date": time.Now().Format("2006-01-02 15:04:05"),
	}
	excel.AddSummary(summary)
	excel.AutoFitColumns(7)

	data, err := excel.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("users_report_%s.xlsx", time.Now().Format("20060102_150405"))
	return data, filename, nil
}

func (s *reportService) GenerateProductsReport(format string) ([]byte, string, error) {
	products, _, err := s.productRepo.List(dto.ProductQuery{
		Page:  1,
		Limit: 1000,
	})
	if err != nil {
		return nil, "", err
	}

	if format == "pdf" {
		return s.generateProductsPDF(products)
	}
	return s.generateProductsExcel(products)
}

func (s *reportService) generateProductsPDF(products []domain.Product) ([]byte, string, error) {
	pdf := utils.NewPDFGenerator()
	pdf.SetTitle("Products Report")

	headers := []string{"No", "Name", "SKU", "Category", "Price", "Stock"}
	widths := []float64{10, 60, 40, 35, 25, 20}
	pdf.AddTableHeader(headers, widths)

	totalValue := 0.0
	totalStock := 0

	for i, product := range products {
		price := product.Price
		stock := product.Stock

		totalValue += price * float64(stock)
		totalStock += stock

		categoryName := ""
		if product.Category != nil {
			categoryName = product.Category.Name
		}

		values := []string{
			fmt.Sprintf("%d", i+1),
			product.Name,
			product.SKU,
			categoryName,
			fmt.Sprintf("Rp %.0f", price),
			fmt.Sprintf("%d", stock),
		}
		pdf.AddTableRow(values, widths)
	}

	summary := map[string]string{
		"Total Products":        fmt.Sprintf("%d", len(products)),
		"Total Stock":           fmt.Sprintf("%d", totalStock),
		"Total Inventory Value": fmt.Sprintf("Rp %.0f", totalValue),
		"Report Date":           time.Now().Format("2006-01-02"),
	}
	pdf.AddSummary("Summary", summary)

	data, err := pdf.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("products_report_%s.pdf", time.Now().Format("20060102_150405"))
	return data, filename, nil
}

func (s *reportService) generateProductsExcel(products []domain.Product) ([]byte, string, error) {
	excel := utils.NewExcelGenerator("Products Report")
	excel.SetTitle("Products Report - " + time.Now().Format("2006-01-02"))

	headers := []string{"No", "Name", "SKU", "Category", "Price", "Stock", "Total Value"}
	excel.AddTableHeader(headers)

	totalValue := 0.0
	totalStock := 0

	for i, product := range products {
		price := product.Price
		stock := product.Stock
		itemValue := price * float64(stock)

		totalValue += itemValue
		totalStock += stock

		categoryName := ""
		if product.Category != nil {
			categoryName = product.Category.Name
		}

		values := []interface{}{
			i + 1,
			product.Name,
			product.SKU,
			categoryName,
			price,
			stock,
			itemValue,
		}
		excel.AddTableRow(values)
	}

	summary := map[string]interface{}{
		"Total Products":        len(products),
		"Total Stock":           totalStock,
		"Total Inventory Value": totalValue,
		"Report Date":           time.Now().Format("2006-01-02 15:04:05"),
	}
	excel.AddSummary(summary)
	excel.AutoFitColumns(7)

	data, err := excel.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("products_report_%s.xlsx", time.Now().Format("20060102_150405"))
	return data, filename, nil
}

func (s *reportService) GenerateOrdersReport(startDate, endDate time.Time, format string) ([]byte, string, error) {
	orders, _, err := s.orderRepo.GetAll(1, 1000)
	if err != nil {
		return nil, "", err
	}

	// Filter by date
	var filteredOrders []domain.Order
	for _, order := range orders {
		if order.CreatedAt.After(startDate) && order.CreatedAt.Before(endDate.Add(24*time.Hour)) {
			filteredOrders = append(filteredOrders, order)
		}
	}

	if format == "pdf" {
		return s.generateOrdersPDF(filteredOrders, startDate, endDate)
	}
	return s.generateOrdersExcel(filteredOrders, startDate, endDate)
}

func (s *reportService) generateOrdersPDF(orders []domain.Order, startDate, endDate time.Time) ([]byte, string, error) {
	pdf := utils.NewPDFGenerator()
	pdf.SetTitle("Orders Report")
	pdf.AddText(fmt.Sprintf("Period: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))

	headers := []string{"No", "Order Number", "Customer", "Amount", "Status", "Date"}
	widths := []float64{10, 45, 45, 30, 25, 35}
	pdf.AddTableHeader(headers, widths)

	totalAmount := 0.0

	for i, order := range orders {
		amount := order.TotalAmount
		totalAmount += amount

		customerName := ""
		if order.User != nil {
			customerName = order.User.Name
		}

		values := []string{
			fmt.Sprintf("%d", i+1),
			order.OrderNumber,
			customerName,
			fmt.Sprintf("Rp %.0f", amount),
			string(order.Status),
			order.CreatedAt.Format("2006-01-02"),
		}
		pdf.AddTableRow(values, widths)
	}

	avgOrder := 0.0
	if len(orders) > 0 {
		avgOrder = totalAmount / float64(len(orders))
	}

	summary := map[string]string{
		"Total Orders":  fmt.Sprintf("%d", len(orders)),
		"Total Revenue": fmt.Sprintf("Rp %.0f", totalAmount),
		"Average Order": fmt.Sprintf("Rp %.0f", avgOrder),
		"Report Date":   time.Now().Format("2006-01-02"),
	}
	pdf.AddSummary("Summary", summary)

	data, err := pdf.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("orders_report_%s.pdf", time.Now().Format("20060102_150405"))
	return data, filename, nil
}

func (s *reportService) generateOrdersExcel(orders []domain.Order, startDate, endDate time.Time) ([]byte, string, error) {
	excel := utils.NewExcelGenerator("Orders Report")
	title := fmt.Sprintf("Orders Report (%s to %s)", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	excel.SetTitle(title)

	headers := []string{"No", "Order Number", "Customer", "Email", "Amount", "Status", "Date"}
	excel.AddTableHeader(headers)

	totalAmount := 0.0

	for i, order := range orders {
		amount := order.TotalAmount
		totalAmount += amount

		customerName := ""
		customerEmail := ""
		if order.User != nil {
			customerName = order.User.Name
			customerEmail = order.User.Email
		}

		values := []interface{}{
			i + 1,
			order.OrderNumber,
			customerName,
			customerEmail,
			amount,
			string(order.Status),
			order.CreatedAt.Format("2006-01-02 15:04"),
		}
		excel.AddTableRow(values)
	}

	avgOrder := 0.0
	if len(orders) > 0 {
		avgOrder = totalAmount / float64(len(orders))
	}

	summary := map[string]interface{}{
		"Total Orders":  len(orders),
		"Total Revenue": totalAmount,
		"Average Order": avgOrder,
		"Report Date":   time.Now().Format("2006-01-02 15:04:05"),
	}
	excel.AddSummary(summary)
	excel.AutoFitColumns(7)

	data, err := excel.Output()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("orders_report_%s.xlsx", time.Now().Format("20060102_150405"))
	return data, filename, nil
}
