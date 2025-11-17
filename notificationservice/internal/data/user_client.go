package data

import (
	"context"

	userv1 "userservice/api/userservice/v1"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

// ProvideBookingClient creates BookingService gRPC client and a cleanup func.
//
// NOTE: change the endpoint to your BookingService gRPC address.
func ProvideUserClient() (userv1.UserServiceClient, func(), error) {
	conn, err := kratosgrpc.DialInsecure(
		context.Background(),
		kratosgrpc.WithEndpoint("127.0.0.1:9000"), // 👈 UserService gRPC address
	)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { _ = conn.Close() }
	return userv1.NewUserServiceClient(conn), cleanup, nil
}

