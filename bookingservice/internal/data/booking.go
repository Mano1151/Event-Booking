package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	v1 "bookingservice/api/bookingservice/v1"
	"bookingservice/internal/biz"

	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type bookingRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

// Booking DB model
type Booking struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	UserID    uint64         `gorm:"index"`
	EventID   uint64         `gorm:"index"`
	SeatIDs   datatypes.JSON `gorm:"type:json"`
	Status    string
	TotalCost float64
	CreatedAt time.Time
}

func NewBookingRepo(db *gorm.DB, redis *redis.Client) biz.BookingRepo {
	db.AutoMigrate(&Booking{})
	return &bookingRepo{db: db, redis: redis}
}

func toProto(b *Booking) *v1.Booking {
	var seatIDs []string
	_ = json.Unmarshal(b.SeatIDs, &seatIDs)
	return &v1.Booking{
		Id:        b.ID,
		UserId:    b.UserID,
		EventId:   b.EventID,
		SeatIds:   seatIDs,
		Status:    b.Status,
		TotalCost: float32(b.TotalCost),
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
	}
}

// ---------------- CRUD ----------------
func (r *bookingRepo) Create(ctx context.Context, booking *v1.Booking) (*v1.Booking, error) {
	seatJSON, _ := json.Marshal(booking.SeatIds)
	b := &Booking{
		UserID:    booking.UserId,
		EventID:   booking.EventId,
		SeatIDs:   seatJSON,
		Status:    booking.Status,
		TotalCost: float64(booking.TotalCost),
		CreatedAt: time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(b).Error; err != nil {
		return nil, err
	}
	return toProto(b), nil
}

func (r *bookingRepo) Get(ctx context.Context, id uint64) (*v1.Booking, error) {
	var b Booking
	if err := r.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, err
	}
	return toProto(&b), nil
}

func (r *bookingRepo) List(ctx context.Context) ([]*v1.Booking, error) {
	var bookings []Booking
	if err := r.db.WithContext(ctx).Find(&bookings).Error; err != nil {
		return nil, err
	}
	res := make([]*v1.Booking, 0, len(bookings))
	for _, b := range bookings {
		res = append(res, toProto(&b))
	}
	return res, nil
}

func (r *bookingRepo) Update(ctx context.Context, booking *v1.Booking) (*v1.Booking, error) {
	seatJSON, _ := json.Marshal(booking.SeatIds)
	var b Booking
	if err := r.db.WithContext(ctx).First(&b, booking.Id).Error; err != nil {
		return nil, err
	}
	b.SeatIDs = seatJSON
	b.Status = booking.Status
	b.TotalCost = float64(booking.TotalCost)
	if err := r.db.WithContext(ctx).Save(&b).Error; err != nil {
		return nil, err
	}
	return toProto(&b), nil
}

func (r *bookingRepo) Cancel(ctx context.Context, id uint64) (*v1.Booking, error) {
	var b Booking
	if err := r.db.WithContext(ctx).First(&b, id).Error; err != nil {
		return nil, err
	}
	b.Status = "CANCELLED"
	if err := r.db.WithContext(ctx).Save(&b).Error; err != nil {
		return nil, err
	}
	return toProto(&b), nil
}

// ---------------- Lock / Unlock ----------------
func (r *bookingRepo) LockSeat(ctx context.Context, eventID uint64, seatID string, userID uint64) (bool, error) {
	key := fmt.Sprintf("booking:lock:%d:%s", eventID, seatID)
	return r.redis.SetNX(ctx, key, userID, 2*time.Minute).Result()
}

func (r *bookingRepo) UnlockSeat(ctx context.Context, eventID uint64, seatID string, userID uint64) error {
	key := fmt.Sprintf("booking:lock:%d:%s", eventID, seatID)
	return r.redis.Del(ctx, key).Err()
}

func (r *bookingRepo) GetLockedSeats(ctx context.Context, eventID uint64) ([]string, error) {
	pattern := fmt.Sprintf("booking:lock:%d:*", eventID)
	keys, err := r.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	seatIDs := []string{}
	for _, k := range keys {
		parts := strings.Split(k, ":")
		if len(parts) == 4 {
			seatIDs = append(seatIDs, parts[3])
		}
	}
	return seatIDs, nil
}

func (r *bookingRepo) ListByEventAndStatus(ctx context.Context, eventID uint64, status string) ([]*v1.Booking, error) {
	var bookings []Booking
	if err := r.db.WithContext(ctx).Where("event_id = ? AND status = ?", eventID, status).Find(&bookings).Error; err != nil {
		return nil, err
	}
	res := make([]*v1.Booking, 0, len(bookings))
	for _, b := range bookings {
		res = append(res, toProto(&b))
	}
	return res, nil
}
