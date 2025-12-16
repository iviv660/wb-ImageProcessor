package image

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
	"github.com/iviv660/wb-ImageProcessor/internal/service"
)

type imageJob struct {
	ImageID int64 `json:"image_id"`
}

func (s *Service) Upload(ctx context.Context, in service.UploadInput) (model.Image, error) {

	originalURL, err := s.storage.Save(ctx, in.FileName, in.Data)
	if err != nil {
		return model.Image{}, fmt.Errorf("storage save: %w", err)
	}

	img, err := s.repo.Create(ctx, model.Image{
		FileName:         in.FileName,
		Format:           in.Format,
		OriginalImageURL: originalURL,
		ChangedImageURL:  nil,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, originalURL)
		return model.Image{}, fmt.Errorf("db create: %w", err)
	}

	val, err := json.Marshal(imageJob{ImageID: img.ID})
	if err != nil {
		_ = s.repo.Delete(ctx, img.ID)
		_ = s.storage.Delete(ctx, originalURL)
		return model.Image{}, fmt.Errorf("marshal job: %w", err)
	}

	key := []byte(strconv.FormatInt(img.ID, 10))
	if err := s.producer.Send(ctx, key, val); err != nil {
		_ = s.repo.Delete(ctx, img.ID) // чтобы не осталось "битой" записи
		_ = s.storage.Delete(ctx, originalURL)
		return model.Image{}, fmt.Errorf("kafka send: %w", err)
	}

	return img, nil
}
