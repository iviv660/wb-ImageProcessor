package image

import (
	"github.com/iviv660/wb-ImageProcessor/internal/adapter"
	"github.com/iviv660/wb-ImageProcessor/internal/repository"
	"github.com/iviv660/wb-ImageProcessor/internal/storage"
)

type Service struct {
	repo     repository.ImageRepository
	storage  storage.Storage
	producer adapter.Producer
}

func New(repo repository.ImageRepository, storage storage.Storage, producer adapter.Producer) *Service {
	return &Service{
		repo:     repo,
		storage:  storage,
		producer: producer,
	}
}
