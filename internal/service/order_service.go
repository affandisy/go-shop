package service

import (
	"fmt"
	"time"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
)

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{orderRepo: orderRepo, productRepo: productRepo}
}

func (s *orderService) CreateOrder(userID string, req dto.CreateOrderRequest) (*domain.Order, error) {
	if len(req.Items) == 0 {
		return nil, domain.ErrEmptyCart
	}

	orderNumber := generateOrderNumber()

	order := &domain.Order{
		OrderNumber: orderNumber,
		UserID:      userID,
		Status:      domain.OrderStatusPending,
		Notes:       req.Notes,
		OrderItems:  []domain.OrderItem{},
	}

	totalAmount := 0.0

	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		if !product.IsAvailable() {
			return nil, fmt.Errorf("product %s is not available", product.Name)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		orderItem := domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		order.OrderItems = append(order.OrderItems, orderItem)
		totalAmount += product.Price * float64(item.Quantity)

		if err := product.ReduceStock(item.Quantity); err != nil {
			return nil, err
		}

		if err := s.productRepo.Update(product); err != nil {
			return nil, err
		}
	}

	order.TotalAmount = totalAmount

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	order, _ = s.orderRepo.GetByID(order.ID)

	return order, nil
}

func (s *orderService) GetOrderByID(orderID string, userID string, isAdmin bool) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if !isAdmin && order.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return order, nil
}

func (s *orderService) GetMyOrders(userID string, page, limit int) ([]domain.Order, int64, error) {
	return s.orderRepo.GetByUserID(userID, page, limit)
}

func (s *orderService) GetAllOrders(page, limit int) ([]domain.Order, int64, error) {
	return s.orderRepo.GetAll(page, limit)
}

func (s *orderService) UpdateOrderStatus(orderID string, req dto.UpdateOrderStatusRequest) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	order.Status = req.Status

	if req.Status == domain.OrderStatusPaid {
		order.MarkAsPaid()
	}

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) CancelOrder(orderID string, userID string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil
	}

	if order.UserID != userID {
		return domain.ErrForbidden
	}

	if !order.CanBeCancelled() {
		return domain.ErrCannotCancelOrder
	}

	order.Status = domain.OrderStatusCancelled

	for _, item := range order.OrderItems {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			continue
		}

		product.Stock += item.Quantity
		s.productRepo.Update(product)
	}

	return s.orderRepo.Update(order)
}

func generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%d", now.Format("20060102"), now.Unix())
}
