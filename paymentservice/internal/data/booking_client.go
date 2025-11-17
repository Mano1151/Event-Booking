package data

import (
	"context"
	"paymentservice/internal/biz"

	bookingv1 "bookingservice/api/bookingservice/v1"
)

type bookingClient struct {
	client bookingv1.BookingServiceClient
}


func (b *bookingClient) UpdateBookingStatus(ctx context.Context, bookingID uint64, status string) error {
	_, err := b.client.UpdateBooking(ctx, &bookingv1.UpdateBookingRequest{
		Id:     bookingID,
		Status: status,
	})
	return err
}

func (b *bookingClient) GetBooking(ctx context.Context, bookingID uint64) (*biz.Booking, error) {
    res, err := b.client.GetBooking(ctx, &bookingv1.GetBookingRequest{Id: bookingID})
    if err != nil {
        return nil, err
    }
    return &biz.Booking{
        ID:        res.Booking.Id,
        UserID:    res.Booking.UserId,
        TotalCost: float64(res.Booking.TotalCost),
        Status:    res.Booking.Status,
    }, nil
}
