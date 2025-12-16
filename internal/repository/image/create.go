package image

import (
	"context"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

func (r *Repository) Create(ctx context.Context, image model.Image) (model.Image, error) {
	query := `
INSERT INTO images(file_name, format, original_url, changed_url) 
VALUES ($1, $2, $3, $4)
RETURNING id, file_name, format, original_url, changed_url;`

	if err := r.pool.QueryRow(ctx, query,
		image.FileName,
		image.Format,
		image.OriginalImageURL,
		image.ChangedImageURL).Scan(
		&image.ID,
		&image.FileName,
		&image.Format,
		&image.OriginalImageURL,
		&image.ChangedImageURL); err != nil {
		return model.Image{}, err
	}

	return image, nil
}

//CREATE TABLE images (
//id BIGSERIAL PRIMARY KEY,
//file_name TEXT NOT NULL,
//format TEXT NOT NULL,
//original_url TEXT NOT NULL,
//changed_url TEXT,
//created_at TIMESTAMPTZ NOT NULL DEFAULT now()
//);
