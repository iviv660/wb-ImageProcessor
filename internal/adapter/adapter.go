package adapter

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}

type Consumer interface {
	Consume(ctx context.Context, handle func(ctx context.Context, msg kafka.Message) error) error
}
