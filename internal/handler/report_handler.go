package handler

import (
	"time"

	"github.com/affandisy/goshop/internal/service"
	"github.com/affandisy/goshop/pkg/response"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GenerateUsersReport(c *gin.Context) {
	format := c.DefaultQuery("format", "pdf")

	if format != "pdf" && format != "excel" {
		response.BadRequest(c, "Invalid format. Use 'pdf' or 'excel'", nil)
		return
	}

	data, filename, err := h.reportService.GenerateUsersReport(format)
	if err != nil {
		response.InternalServerError(c, "Failed to generate report", nil)
		return
	}

	contentType := "application/pdf"
	if format == "excel" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, contentType, data)
}

func (h *ReportHandler) GenerateProductsReport(c *gin.Context) {
	format := c.DefaultQuery("format", "pdf")

	if format != "pdf" && format != "excel" {
		response.BadRequest(c, "Invalid format. Use 'pdf' or 'excel'", nil)
		return
	}

	data, filename, err := h.reportService.GenerateProductsReport(format)
	if err != nil {
		response.InternalServerError(c, "Failed to generate report", nil)
		return
	}

	contentType := "application/pdf"
	if format == "excel" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, contentType, data)
}

func (h *ReportHandler) GenerateOrdersReport(c *gin.Context) {
	format := c.DefaultQuery("format", "pdf")

	if format != "pdf" && format != "excel" {
		response.BadRequest(c, "Invalid format. Use 'pdf' or 'excel'", nil)
		return
	}

	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		response.BadRequest(c, "Invalid start_date format. Use YYYY-MM-DD", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		response.BadRequest(c, "Invalid end_date format. Use YYYY-MM-DD", err)
		return
	}

	data, filename, err := h.reportService.GenerateOrdersReport(startDate, endDate, format)
	if err != nil {
		response.InternalServerError(c, "Failed to generate report", nil)
		return
	}

	contentType := "application/pdf"
	if format == "excel" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, contentType, data)
}

func (h *ReportHandler) GetReportTypes(c *gin.Context) {
	reportTypes := []map[string]string{
		{
			"name":        "Users Report",
			"endpoint":    "/api/v1/reports/users",
			"description": "Report of all users in the system",
			"formats":     "pdf, excel",
		},
		{
			"name":        "Products Report",
			"endpoint":    "/api/v1/reports/products",
			"description": "Report of all products with inventory value",
			"formats":     "pdf, excel",
		},
		{
			"name":        "Orders Report",
			"endpoint":    "/api/v1/reports/orders",
			"description": "Report of orders with date range filter",
			"formats":     "pdf, excel",
			"parameters":  "start_date, end_date",
		},
	}

	response.Success(c, "Available reports", reportTypes)
}
