package data

import (
	"context"

	bookingv1 "bookingservice/api/bookingservice/v1"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

// ProvideBookingClient creates BookingService gRPC client and a cleanup func.
//
// NOTE: change the endpoint to your BookingService gRPC address.
func ProvideBookingClient() (bookingv1.BookingServiceClient, func(), error) {
	conn, err := kratosgrpc.DialInsecure(
		context.Background(),
		kratosgrpc.WithEndpoint("127.0.0.1:9002"), // <-- set BookingService gRPC address/port
	)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = conn.Close() }
	return bookingv1.NewBookingServiceClient(conn), cleanup, nil
}
