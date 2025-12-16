package image

import (
	"context"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

func (r *Repository) GetByID(ctx context.Context, id int64) (model.Image, error) {
	query := `
SELECT id, file_name, format, original_url, changed_url
FROM images
WHERE id = $1`

	var image model.Image

	if err := r.pool.QueryRow(ctx, query, id).Scan(
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
