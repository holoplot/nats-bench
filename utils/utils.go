package utils

import (
	"context"
	"fmt"

	nats "github.com/nats-io/nats.go"
)

type OnConnected func(ctx context.Context, nc *nats.Conn)

func NatsConnect(ctx context.Context, natsURL, clientID string, onConnected OnConnected) error {
	opts := []nats.Option{
		nats.Name(clientID),
		nats.DisconnectErrHandler(
			func(_ *nats.Conn, err error) {
				panic(err)
			},
		),
		nats.ConnectHandler(
			func(nc *nats.Conn) {
				onConnected(ctx, nc)
			},
		),
	}

	_, err := nats.Connect(natsURL, opts...)
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	return nil
}
