package image

import (
	"context"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM images WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
