package repository

import (
	"context"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

type ImageRepository interface {
	Create(ctx context.Context, image model.Image) (model.Image, error)
	GetByID(ctx context.Context, id int64) (model.Image, error)
	Update(ctx context.Context, image model.Image) (model.Image, error)
	Delete(ctx context.Context, id int64) error
}
