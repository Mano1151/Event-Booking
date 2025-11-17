package data

import (
	"context"
    "fmt"
	v1 "userservice/api/userservice/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// ProvideUserClient returns a UserServiceClient and cleanup
func ProvideUserClient() (v1.UserServiceClient, func(), error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"), // UserService gRPC port
	)
	if err != nil {
    return nil, nil, fmt.Errorf("failed to connect to user service: %w", err)
}

	cleanup := func() {
		_ = conn.Close()
	}

	client := v1.NewUserServiceClient(conn)
	return client, cleanup, nil
}
