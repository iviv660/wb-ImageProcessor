package image

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (l *LocalStorage) Save(ctx context.Context, filename string, r io.Reader) (string, error) {
	if err := os.MkdirAll(l.basePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	dstName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(filename))
	dstPath := filepath.Join(l.basePath, dstName)

	f, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create dst file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	url := l.baseURL + "/" + dstName
	return url, nil
}
