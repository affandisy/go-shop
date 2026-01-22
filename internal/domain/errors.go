package domain

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")

	// Product errors
	ErrProductNotFound     = errors.New("product not found")
	ErrProductNotAvailable = errors.New("product not available")
	ErrInsufficientStock   = errors.New("insufficient stock")

	// Category errors
	ErrCategoryNotFound = errors.New("category not found")

	// Order errors
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrCannotCancelOrder  = errors.New("cannot cancel order")
	ErrEmptyCart          = errors.New("cart is empty")

	// Payment errors
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentAlreadyExists = errors.New("payment already exists for this order")
	ErrPaymentFailed        = errors.New("payment failed")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
	ErrOrderAlreadyPaid     = errors.New("order already paid")

	// General errors
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternalServer = errors.New("internal server error")
	ErrForbidden      = errors.New("forbidden")
)
