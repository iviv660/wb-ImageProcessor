package image

import (
	"context"
	"fmt"
)

func (s *Service) DeleteImage(ctx context.Context, id int64) error {
	img, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get image: %w", err)
	}

	if img.OriginalImageURL != "" {
		if err := s.storage.Delete(ctx, img.OriginalImageURL); err != nil {
			return fmt.Errorf("delete original: %w", err)
		}
	}
	if img.ChangedImageURL != nil {
		if err := s.storage.Delete(ctx, *img.ChangedImageURL); err != nil {
			return fmt.Errorf("delete changed: %w", err)
		}
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete db: %w", err)
	}

	return nil
}
