package biz

import (
	"context"
	"fmt"

	bookingv1 "bookingservice/api/bookingservice/v1"
	eventv1 "eventservice/api/eventservice/v1"

	"github.com/go-kratos/kratos/v2/log"
	
)

type BookingRepo interface {
	Create(ctx context.Context, booking *bookingv1.Booking) (*bookingv1.Booking, error)
	Get(ctx context.Context, id uint64) (*bookingv1.Booking, error)
	List(ctx context.Context) ([]*bookingv1.Booking, error)
	Update(ctx context.Context, booking *bookingv1.Booking) (*bookingv1.Booking, error)
	Cancel(ctx context.Context, id uint64) (*bookingv1.Booking, error)
	LockSeat(ctx context.Context, eventID uint64, seatID string, userID uint64) (bool, error)
	UnlockSeat(ctx context.Context, eventID uint64, seatID string, userID uint64) error
	GetLockedSeats(ctx context.Context, eventID uint64) ([]string, error)
	ListByEventAndStatus(ctx context.Context, eventID uint64, status string) ([]*bookingv1.Booking, error)
}

type BookingUsecase struct {
	repo        BookingRepo
	eventClient eventv1.EventServiceClient
	log         *log.Helper
}

func NewBookingUsecase(repo BookingRepo, eventClient eventv1.EventServiceClient, logger log.Logger) *BookingUsecase {
	return &BookingUsecase{
		repo:        repo,
		eventClient: eventClient,
		log:         log.NewHelper(logger),
	}
}

// Lock / Unlock
func (uc *BookingUsecase) LockSeat(ctx context.Context, eventID uint64, seatID string, userID uint64) (bool, error) {
	return uc.repo.LockSeat(ctx, eventID, seatID, userID)
}

func (uc *BookingUsecase) UnlockSeat(ctx context.Context, eventID uint64, seatID string, userID uint64) error {
	return uc.repo.UnlockSeat(ctx, eventID, seatID, userID)
}

// Get locked / booked seats
func (uc *BookingUsecase) GetLockedSeats(ctx context.Context, eventID uint64) ([]string, error) {
	uc.log.Infof("Fetching locked seats for event_id=%d", eventID)
	return uc.repo.GetLockedSeats(ctx, eventID)
}

func (uc *BookingUsecase) GetBookedSeats(ctx context.Context, eventID uint64) ([]string, error) {
	uc.log.Infof("Fetching booked seats for event_id=%d", eventID)
	bookings, err := uc.repo.ListByEventAndStatus(ctx, eventID, "CONFIRMED")
	if err != nil {
		return nil, err
	}
	var seatIDs []string
	for _, b := range bookings {
		seatIDs = append(seatIDs, b.SeatIds...)
	}
	return seatIDs, nil
}

// CRUD
func (uc *BookingUsecase) Create(ctx context.Context, req *bookingv1.CreateBookingRequest) (*bookingv1.Booking, error) {
    // 1️⃣ Validate user
    valid, err := uc.eventClient.ValidateUser(ctx, &eventv1.ValidateUserRequest{Id: req.UserId})
    if err != nil || !valid.Found {
        return nil, fmt.Errorf("user not found")
    }

    // 2️⃣ Get event
    evResp, err := uc.eventClient.GetShowEvent(ctx, &eventv1.GetShowEventRequest{Id: req.EventId})
    if err != nil {
        return nil, fmt.Errorf("event not found")
    }
    ev := evResp.ShowEvent

    // 3️⃣ Check available seats
    if int32(len(req.SeatIds)) > ev.AvailableSeats {
        return nil, fmt.Errorf("not enough seats available")
    }

    // 4️⃣ Lock seats
    lockedSeats := []string{}
    for _, seatID := range req.SeatIds {
        locked, _ := uc.repo.LockSeat(ctx, req.EventId, seatID, req.UserId)
        if !locked {
            for _, s := range lockedSeats {
                _ = uc.repo.UnlockSeat(ctx, req.EventId, s, req.UserId)
            }
            return nil, fmt.Errorf("seat already taken")
        }
        lockedSeats = append(lockedSeats, seatID)
    }

    // 5️⃣ Calculate total cost
    totalCost := float32(len(req.SeatIds)) * ev.PricePerSeat

    // 6️⃣ Create booking
    booking := &bookingv1.Booking{
        UserId:    req.UserId,
        EventId:   req.EventId,
        SeatIds:   req.SeatIds,
        Status:    "PENDING",
        TotalCost: totalCost,
    }

    return uc.repo.Create(ctx, booking)
}



func (uc *BookingUsecase) ConfirmBooking(ctx context.Context, bookingID uint64) (*bookingv1.Booking, error) {
	booking, err := uc.repo.Get(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	booking.Status = "CONFIRMED"
	updated, err := uc.repo.Update(ctx, booking)
	if err != nil {
		return nil, err
	}
	if len(booking.SeatIds) > 0 {
		_, _ = uc.eventClient.DecrementSeats(ctx, &eventv1.DecrementSeatsRequest{
			EventId: booking.EventId,
			SeatIds: booking.SeatIds,
		})
	}
	return updated, nil
}

func (uc *BookingUsecase) Cancel(ctx context.Context, bookingID uint64) (*bookingv1.Booking, error) {
	return uc.repo.Cancel(ctx, bookingID)
}

func (uc *BookingUsecase) Get(ctx context.Context, id uint64) (*bookingv1.Booking, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *BookingUsecase) Update(ctx context.Context, booking *bookingv1.Booking) (*bookingv1.Booking, error) {
	return uc.repo.Update(ctx, booking)
}

func (uc *BookingUsecase) List(ctx context.Context) ([]*bookingv1.Booking, error) {
	return uc.repo.List(ctx)
}
