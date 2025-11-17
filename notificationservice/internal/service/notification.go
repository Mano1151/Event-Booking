package service

import (
    "context"
    "fmt"
    userv1 "userservice/api/userservice/v1"
    notifv1 "notificationservice/api/notificationservice/v1"
    "notificationservice/internal/biz"
    bookingv1 "bookingservice/api/bookingservice/v1"

    "github.com/go-kratos/kratos/v2/log"
)

type NotificationService struct {
    notifv1.UnimplementedNotificationServiceServer
    uc            *biz.NotificationUsecase
    bookingClient bookingv1.BookingServiceClient
    userClient    userv1.UserServiceClient
    log           *log.Helper
}

func NewNotificationService(
    uc *biz.NotificationUsecase,
    bookingClient bookingv1.BookingServiceClient,
    userClient userv1.UserServiceClient,
    logger log.Logger,
) *NotificationService {
    return &NotificationService{
        uc:            uc,
        bookingClient: bookingClient,
        userClient:    userClient,
        log:           log.NewHelper(logger),
    }
}

func (s *NotificationService) SendBookingNotification(ctx context.Context, req *notifv1.SendBookingNotificationRequest) (*notifv1.SendBookingNotificationReply, error) {
    // 1️⃣ Fetch booking
    bookingResp, err := s.bookingClient.GetBooking(ctx, &bookingv1.GetBookingRequest{Id: req.BookingId})
    if err != nil || bookingResp.Booking == nil {
        s.log.Errorf("BookingID=%d not found: %v", req.BookingId, err)
        return &notifv1.SendBookingNotificationReply{
            Success: false,
            Message: "Booking not found",
        }, err
    }
    booking := bookingResp.Booking

    // 2️⃣ Only confirmed bookings
    if booking.Status != "CONFIRMED" {
        return &notifv1.SendBookingNotificationReply{
            Success: false,
            Message: "Booking not confirmed",
        }, nil
    }

    // 3️⃣ Fetch user
    userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: booking.UserId})
    if err != nil {
        s.log.Errorf("Failed to fetch user for BookingID=%d, UserID=%d: %v", booking.Id, booking.UserId, err)
        return &notifv1.SendBookingNotificationReply{
            Success: false,
            Message: "Failed to fetch user info",
        }, err
    }

    s.log.Infof("BookingID=%d, UserID=%d, Email=%s", booking.Id, booking.UserId, userResp.Email)

    // 4️⃣ Prepare notification
    body := fmt.Sprintf(
        "Hello %s,\n\nYour booking has been confirmed!\n\nBooking ID: %d\nEvent ID: %d\nSeats: %v\nTotal Cost: %.2f\nBooked At: %s\nStatus: %s",
        userResp.Name,
        booking.Id,
        booking.EventId,
        booking.SeatIds,
        booking.TotalCost,
        booking.CreatedAt,
        booking.Status,
    )

    notif := &biz.Notification{
        BookingID: uint64(booking.Id),
        Email:     userResp.Email,
        Subject:   "Booking Confirmed",
        Body:      body,
        Status:    booking.Status,
    }

    // 5️⃣ Send
    if err := s.uc.Send(ctx, notif); err != nil {
        s.log.Errorf("Failed to send notification for BookingID=%d to %s: %v", booking.Id, userResp.Email, err)
        return &notifv1.SendBookingNotificationReply{
            Success: false,
            Message: "Failed to send notification",
        }, err
    }

    s.log.Infof("Notification sent successfully for BookingID=%d to %s", booking.Id, userResp.Email)

    return &notifv1.SendBookingNotificationReply{
        Success: true,
        Message: "Notification sent successfully",
    }, nil
}
