package data

import (
	"context"
	notifv1 "notificationservice/api/notificationservice/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// ProvideNotificationClient creates a gRPC client to NotificationService
func ProvideNotificationClient() (notifv1.NotificationServiceClient, func(), error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9004"), // NotificationService gRPC port
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	client := notifv1.NewNotificationServiceClient(conn)
	return client, cleanup, nil
}
