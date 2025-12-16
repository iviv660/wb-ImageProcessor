package service

import (
	"context"
	"io"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

type UploadInput struct {
	FileName string
	Format   string
	Size     int64
	Data     io.Reader
}

type ImageService interface {
	Upload(ctx context.Context, in UploadInput) (model.Image, error)
	GetImage(ctx context.Context, id int64) (model.Image, error)
	DeleteImage(ctx context.Context, id int64) error
	ProcessImage(ctx context.Context, id int64) (any, error)
}
