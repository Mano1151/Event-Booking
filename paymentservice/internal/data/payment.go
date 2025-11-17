package data

import (
	"context"
	"time"

	"paymentservice/internal/biz"

	"gorm.io/gorm"
)

// DB model
type PaymentModel struct {
	ID        uint64  `gorm:"primaryKey;autoIncrement"`
	BookingID uint64  `gorm:"not null;index"`
	Amount    float64 `gorm:"not null"`
	Method    string  `gorm:"size:20"`
	Status    string  `gorm:"size:20;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type paymentRepo struct {
	data *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) biz.PaymentRepo {
	return &paymentRepo{data: db}
}

func (r *paymentRepo) Save(ctx context.Context, p *biz.Payment) (*biz.Payment, error) {
	model := &PaymentModel{
		BookingID: p.BookingID,
		Amount:    p.Amount,
		Method:    p.Method,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
	}
	err := r.data.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	p.ID = model.ID
	return p, nil
}

func (r *paymentRepo) UpdateStatus(ctx context.Context, bookingID uint64, status string) error {
	return r.data.WithContext(ctx).
		Model(&PaymentModel{}).
		Where("booking_id = ?", bookingID).
		Update("status", status).Error
}

func (r *paymentRepo) FindByBooking(ctx context.Context, bookingID uint64) (*biz.Payment, error) {
	var model PaymentModel
	err := r.data.WithContext(ctx).Where("booking_id = ?", bookingID).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &biz.Payment{
		ID:        model.ID,
		BookingID: model.BookingID,
		Amount:    model.Amount,
		Method:    model.Method,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
	}, nil
}
