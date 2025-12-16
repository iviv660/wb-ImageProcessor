package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func (p *KafkaProducer) Send(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}
