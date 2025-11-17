package service

import (
	"context"
	v1 "bookingservice/api/bookingservice/v1"
	eventv1 "eventservice/api/eventservice/v1"
	notifv1 "notificationservice/api/notificationservice/v1"
	"bookingservice/internal/biz"
    "time"
	"github.com/go-kratos/kratos/v2/log"
)

type BookingService struct {
	v1.UnimplementedBookingServiceServer
	uc                  *biz.BookingUsecase
	log                 *log.Helper
	eventClient         eventv1.EventServiceClient
	notificationClient  notifv1.NotificationServiceClient
}

func NewBookingService(
	uc *biz.BookingUsecase,
	notificationClient notifv1.NotificationServiceClient,
	eventClient eventv1.EventServiceClient,
	logger log.Logger,
) *BookingService {
	return &BookingService{
		uc:                 uc,
		eventClient:        eventClient,
		notificationClient: notificationClient,
		log:                log.NewHelper(logger),
	}
}

// ------------------- CRUD -------------------
func (s *BookingService) CreateBooking(ctx context.Context, req *v1.CreateBookingRequest) (*v1.CreateBookingReply, error) {
    // 🔍 Log the incoming request to see what UserId is being sent
    s.log.Infof("Incoming CreateBookingRequest: UserId=%d, EventId=%d, SeatIds=%v", req.UserId, req.EventId, req.SeatIds)

    // 1️⃣ Create booking
    booking, err := s.uc.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    s.log.Infof("Booking created: Id=%d, UserId=%d, EventId=%d", booking.Id, booking.UserId, booking.EventId)

    // 2️⃣ Send notification asynchronously
    go func(bookingID uint64) {
        ctxNotif, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        _, err := s.notificationClient.SendBookingNotification(ctxNotif, &notifv1.SendBookingNotificationRequest{
            BookingId: bookingID,
        })
        if err != nil {
            s.log.Errorf("Failed to send booking notification for booking %d: %v", bookingID, err)
        } else {
            s.log.Infof("Notification sent successfully for booking %d", bookingID)
        }
    }(booking.Id)

    return &v1.CreateBookingReply{Booking: booking}, nil
}


func (s *BookingService) GetBooking(ctx context.Context, req *v1.GetBookingRequest) (*v1.CreateBookingReply, error) {
	booking, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.CreateBookingReply{Booking: booking}, nil
}

func (s *BookingService) ListBookings(ctx context.Context, req *v1.ListBookingsRequest) (*v1.ListBookingsReply, error) {
	bookings, err := s.uc.List(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ListBookingsReply{Bookings: bookings}, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *v1.UpdateBookingRequest) (*v1.UpdateBookingReply, error) {
	// 1️⃣ Fetch existing booking
	booking, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return &v1.UpdateBookingReply{Success: false}, err
	}

	// 2️⃣ Update status
	booking.Status = req.Status
	updatedBooking, err := s.uc.Update(ctx, booking)
	if err != nil {
		return &v1.UpdateBookingReply{Success: false}, err
	}

	eventID := updatedBooking.EventId
	seatIDs := updatedBooking.SeatIds

	switch updatedBooking.Status {
	case "CONFIRMED":
		// a) Decrement seats
		_, err := s.eventClient.DecrementSeats(ctx, &eventv1.DecrementSeatsRequest{
			EventId: eventID,
			SeatIds: seatIDs,
		})
		if err != nil {
			s.log.Errorf("Failed to decrement seats: %v", err)
			return &v1.UpdateBookingReply{Success: false}, err
		}

		// b) Send confirmation notification
		go func(bookingID uint64) {
			ctxNotif, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := s.notificationClient.SendBookingNotification(ctxNotif, &notifv1.SendBookingNotificationRequest{
				BookingId: bookingID, // notification service fetches email
			})
			if err != nil {
				s.log.Errorf("Failed to send confirmation notification for booking %d: %v", bookingID, err)
			} else {
				s.log.Infof("Confirmation notification sent for booking %d", bookingID)
			}
		}(updatedBooking.Id)

	case "CANCELLED":
		// a) Increment seats
		_, err := s.eventClient.IncrementSeats(ctx, &eventv1.IncrementSeatsRequest{
			EventId: eventID,
			SeatIds: seatIDs,
		})
		if err != nil {
			s.log.Errorf("Failed to increment seats: %v", err)
			return &v1.UpdateBookingReply{Success: false}, err
		}

		// b) Send cancellation notification
		go func(bookingID uint64) {
			ctxNotif, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := s.notificationClient.SendBookingNotification(ctxNotif, &notifv1.SendBookingNotificationRequest{
				BookingId: bookingID,
			})
			if err != nil {
				s.log.Errorf("Failed to send cancellation notification for booking %d: %v", bookingID, err)
			} else {
				s.log.Infof("Cancellation notification sent for booking %d", bookingID)
			}
		}(updatedBooking.Id)
	}

	s.log.Infof("Booking updated: Id=%d, Status=%s", updatedBooking.Id, updatedBooking.Status)
	return &v1.UpdateBookingReply{Success: true}, nil
}


func (s *BookingService) CancelBooking(ctx context.Context, req *v1.CancelBookingRequest) (*v1.CreateBookingReply, error) {
	booking, err := s.uc.Cancel(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.CreateBookingReply{Booking: booking}, nil
}

func (s *BookingService) ConfirmBooking(ctx context.Context, req *v1.ConfirmBookingRequest) (*v1.CreateBookingReply, error) {
	booking, err := s.uc.ConfirmBooking(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.CreateBookingReply{Booking: booking}, nil
}

// ------------------- Seats -------------------
func (s *BookingService) GetBookedSeats(ctx context.Context, req *v1.GetBookedSeatsRequest) (*v1.GetBookedSeatsReply, error) {
	bookedSeats, err := s.uc.GetBookedSeats(ctx, req.EventId)
	if err != nil {
		return &v1.GetBookedSeatsReply{SeatIds: []string{}}, err
	}
	return &v1.GetBookedSeatsReply{SeatIds: bookedSeats}, nil
}

func (s *BookingService) GetLockedSeats(ctx context.Context, req *v1.GetLockedSeatsRequest) (*v1.GetLockedSeatsReply, error) {
	lockedSeats, err := s.uc.GetLockedSeats(ctx, req.EventId)
	if err != nil {
		return &v1.GetLockedSeatsReply{SeatIds: []string{}}, err
	}
	return &v1.GetLockedSeatsReply{SeatIds: lockedSeats}, nil
}

func (s *BookingService) LockSeat(ctx context.Context, req *v1.LockSeatRequest) (*v1.LockSeatReply, error) {
	var allLocked bool = true
	for _, seatId := range req.SeatIds {
		locked, err := s.uc.LockSeat(ctx, req.EventId, seatId, req.UserId)
		if err != nil || !locked {
			allLocked = false
			// optional: log the seat that failed
			continue
		}
	}
	return &v1.LockSeatReply{Locked: allLocked}, nil
}

func (s *BookingService) UnlockSeat(ctx context.Context, req *v1.UnlockSeatRequest) (*v1.UnlockSeatReply, error) {
	var allSuccess bool = true
	for _, seatId := range req.SeatIds {
		err := s.uc.UnlockSeat(ctx, req.EventId, seatId, req.UserId)
		if err != nil {
			allSuccess = false
			// optional: log the seat that failed
			continue
		}
	}
	return &v1.UnlockSeatReply{Success: allSuccess}, nil
}

