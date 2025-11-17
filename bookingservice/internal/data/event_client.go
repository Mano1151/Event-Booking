package data

import (
	"context"
	eventv1 "eventservice/api/eventservice/v1"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// ProvideEventClient creates a gRPC client to EventService
func ProvideEventClient() (eventv1.EventServiceClient, func(), error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9001"), // EventService gRPC port
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	client := eventv1.NewEventServiceClient(conn)
	return client, cleanup, nil
}
