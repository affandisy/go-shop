package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/affandisy/goshop/internal/domain"
	"github.com/affandisy/goshop/internal/domain/dto"
	"github.com/affandisy/goshop/internal/repository"
	"github.com/affandisy/goshop/pkg/payment"
)

type paymentService struct {
	paymentRepo    repository.PaymentRepository
	orderRepo      repository.OrderRepository
	midtransClient *payment.MidtransClient
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, midtransClient *payment.MidtransClient) PaymentService {
	return &paymentService{paymentRepo: paymentRepo, orderRepo: orderRepo, midtransClient: midtransClient}
}

func (s *paymentService) CreatePayment(userID string, req dto.CreatePaymentRequest) (*domain.Payment, error) {
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if order.Status != domain.OrderStatusPending {
		return nil, domain.ErrOrderAlreadyPaid
	}

	existingPayment, err := s.paymentRepo.GetByOrderID(req.OrderID)
	if err != nil && !errors.Is(err, domain.ErrPaymentNotFound) {
		return nil, err
	}

	if existingPayment != nil && existingPayment.Status == domain.PaymentStatusSuccess {
		return nil, domain.ErrPaymentAlreadyExists
	}

	midtransOrderID := fmt.Sprintf("PAY-%s-%d", order.OrderNumber, time.Now().Unix())

	snapReq := payment.CreateSnapTokenRequest{
		OrderID:       midtransOrderID,
		GrossAmount:   int64(order.TotalAmount),
		CustomerName:  order.User.Name,
		CustomerEmail: order.User.Email,
		CustomerPhone: order.User.Phone,
		Items:         []payment.ItemDetail{},
	}

	for _, item := range order.OrderItems {
		snapReq.Items = append(snapReq.Items, payment.ItemDetail{
			ID:       item.ProductID,
			Name:     item.Product.Name,
			Price:    int64(item.Price),
			Quantity: int32(item.Quantity),
		})
	}

	snapResp, err := s.midtransClient.CreateSnapToken(snapReq)
	if err != nil {
		return nil, err
	}

	var paymentRecord *domain.Payment

	if existingPayment != nil {
		existingPayment.MidtransOrderID = midtransOrderID
		existingPayment.MidtransSnapToken = snapResp.Token
		existingPayment.MidtransSnapURL = snapResp.RedirectURL
		existingPayment.Status = domain.PaymentStatusPending
		existingPayment.PaymentMethod = req.PaymentMethod

		expiredAt := time.Now().Add(24 * time.Hour)
		existingPayment.ExpiredAt = &expiredAt

		if err := s.paymentRepo.Update(existingPayment); err != nil {
			return nil, err
		}

		paymentRecord = existingPayment
	} else {
		expiredAt := time.Now().Add(24 * time.Hour)

		paymentRecord = &domain.Payment{
			OrderID:           order.ID,
			Amount:            order.TotalAmount,
			PaymentMethod:     req.PaymentMethod,
			Status:            domain.PaymentStatusPending,
			MidtransOrderID:   midtransOrderID,
			MidtransSnapToken: snapResp.Token,
			MidtransSnapURL:   snapResp.RedirectURL,
			ExpiredAt:         &expiredAt,
		}

		if err := s.paymentRepo.Create(paymentRecord); err != nil {
			return nil, err
		}
	}

	paymentRecord, _ = s.paymentRepo.GetByID(paymentRecord.ID)

	return paymentRecord, nil
}

func (s *paymentService) GetPaymentByID(id string) (*domain.Payment, error) {
	return s.paymentRepo.GetByID(id)
}

func (s *paymentService) GetPaymentByOrderID(orderID string) (*domain.Payment, error) {
	return s.paymentRepo.GetByOrderID(orderID)
}
func (s *paymentService) HandleNotification(notification dto.PaymentNotification) error {
	payment, err := s.paymentRepo.GetByMidtransOrderID(notification.OrderID)
	if err != nil {
		return err
	}

	order := payment.Order

	switch notification.TransactionStatus {
	case "capture", "settlement":
		if notification.FraudStatus == "accept" || notification.FraudStatus == "" {
			payment.MarkAsPaid()
			payment.PaymentMethod = domain.PaymentMethod(notification.PaymentType)
			order.MarkAsPaid()
		}
	case "pending":
		payment.Status = domain.PaymentStatusPending
	case "deny", "cancel":
		payment.MarkAsFailed()
	case "expire":
		payment.MarkAsExpired()
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	if order.Status == domain.OrderStatusPaid {
		if err := s.orderRepo.Update(order); err != nil {
			return err
		}
	}

	return nil
}

func (s *paymentService) GetAllPayments(page, limit int) ([]domain.Payment, int64, error) {
	return s.paymentRepo.List(page, limit)
}
