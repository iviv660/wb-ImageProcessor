package consumer

import "github.com/segmentio/kafka-go"

type KafkaConsumer struct {
	reader *kafka.Reader
}

func New(reader *kafka.Reader) *KafkaConsumer {
	return &KafkaConsumer{
		reader: reader,
	}
}
