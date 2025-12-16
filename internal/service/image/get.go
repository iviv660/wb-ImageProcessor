package image

import (
	"context"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

func (s *Service) GetImage(ctx context.Context, id int64) (model.Image, error) {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Image{}, err
	}
	return image, nil
}
