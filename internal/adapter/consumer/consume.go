package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func (c *KafkaConsumer) Consume(ctx context.Context, handle func(ctx context.Context, msg kafka.Message) error) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		if err := handle(ctx, msg); err != nil {
			return err
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}
