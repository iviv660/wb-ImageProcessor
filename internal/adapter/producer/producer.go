package producer

import "github.com/segmentio/kafka-go"

type KafkaProducer struct {
	writer *kafka.Writer
}

func New(writer *kafka.Writer) *KafkaProducer {
	return &KafkaProducer{
		writer: writer,
	}
}
