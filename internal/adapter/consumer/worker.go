package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iviv660/wb-ImageProcessor/internal/adapter"
	"github.com/segmentio/kafka-go"
)

type ImageService interface {
	ProcessImage(ctx context.Context, id int64) (any, error)
}

type Worker struct {
	consumer adapter.Consumer
	svc      ImageService
}

func NewWorker(c adapter.Consumer, svc ImageService) *Worker {
	return &Worker{
		consumer: c,
		svc:      svc,
	}
}

type imageJob struct {
	ImageID int64 `json:"image_id"`
}

func (w *Worker) Run(ctx context.Context) error {
	return w.consumer.Consume(ctx, func(ctx context.Context, msg kafka.Message) error {
		var job imageJob
		if err := json.Unmarshal(msg.Value, &job); err != nil {
			return fmt.Errorf("bad message: %w", err)
		}
		if job.ImageID <= 0 {
			return fmt.Errorf("invalid image_id: %d", job.ImageID)
		}
		_, err := w.svc.ProcessImage(ctx, job.ImageID)
		if err != nil {
			return fmt.Errorf("process image %d: %w", job.ImageID, err)
		}
		return nil
	})
}
