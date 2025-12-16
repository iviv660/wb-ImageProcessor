package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (l *LocalStorage) Delete(ctx context.Context, url string) error {
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return fmt.Errorf("invalid URL: %s", url)
	}
	name := parts[len(parts)-1]
	path := filepath.Join(l.basePath, name)
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}
	return nil
}
