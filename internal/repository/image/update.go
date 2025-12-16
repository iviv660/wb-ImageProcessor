package image

import (
	"context"
	"errors"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Update(ctx context.Context, image model.Image) (model.Image, error) {
	query := `
UPDATE images
SET file_name = $2,
    format = $3,
    original_url = $4,
    changed_url = $5
WHERE id = $1
RETURNING id, file_name, format, original_url, changed_url;`

	err := r.pool.QueryRow(ctx, query,
		image.ID,
		image.FileName,
		image.Format,
		image.OriginalImageURL,
		image.ChangedImageURL,
	).Scan(
		&image.ID,
		&image.FileName,
		&image.Format,
		&image.OriginalImageURL,
		&image.ChangedImageURL,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Image{}, model.ErrNotFound
		}
		return model.Image{}, err
	}

	return image, nil
}
