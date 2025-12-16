package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, filename string, r io.Reader) (string, error)
	Open(ctx context.Context, url string) (io.ReadCloser, error)
	Delete(ctx context.Context, url string) error
}
