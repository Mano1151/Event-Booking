package biz

import (
	"context"
	"fmt"
	"time"
	"math/rand"

	
	
)

// Payment entity
type Payment struct {
	ID        uint64
	BookingID uint64
	Amount    float64
	Method    string
	Status    string
	CreatedAt time.Time
}


// PaymentRepo interface
type PaymentRepo interface {
	Save(ctx context.Context, p *Payment) (*Payment, error)
	UpdateStatus(ctx context.Context, bookingID uint64, status string) error
	FindByBooking(ctx context.Context, bookingID uint64) (*Payment, error)
}

// BookingClient interface to fetch booking info
type BookingClient interface {
	GetBooking(ctx context.Context, bookingID uint64) (*Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID uint64, status string) error
}

// Booking struct (minimal fields for payment)
type Booking struct {
	ID        uint64
	UserID    uint64
	TotalCost float64
	Status    string
}

// PaymentUsecase
type PaymentUsecase struct {
	repo          PaymentRepo
	bookingClient BookingClient
}

func NewPaymentUsecase(repo PaymentRepo, bc BookingClient) *PaymentUsecase {
	return &PaymentUsecase{
		repo:          repo,
		bookingClient: bc,
	}
}

// ProcessPayment creates payment in DB and updates status
func (uc *PaymentUsecase) ProcessPayment(ctx context.Context, bookingID uint64, method string) (*Payment, error) {
	// 1️⃣ Fetch booking from BookingService
	booking, err := uc.bookingClient.GetBooking(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}

	// 2️⃣ Create payment record with PENDING status
	payment := &Payment{
		BookingID: bookingID,
		Amount:    booking.TotalCost, // amount comes from booking
		Method:    method,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	savedPayment, err := uc.repo.Save(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// 3️⃣ Randomly determine payment success/failure
	success := rand.Intn(2) == 0 // 50% chance

	if success {
		savedPayment.Status = "PAID"
		booking.Status = "CONFIRMED"
	} else {
		savedPayment.Status = "FAILED"
		booking.Status = "CANCELLED"
	}

	// 4️⃣ Update payment status in database
	if err := uc.repo.UpdateStatus(ctx, bookingID, savedPayment.Status); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// 5️⃣ Update booking status in BookingService via gRPC
	err = uc.bookingClient.UpdateBookingStatus(ctx, bookingID, booking.Status)

	if err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	return savedPayment, nil
}