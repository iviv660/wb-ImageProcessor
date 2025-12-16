package image

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func (l *LocalStorage) Open(ctx context.Context, urlStr string) (io.ReadCloser, error) {
	u, err := url.Parse(urlStr)
	if err == nil && u.Path != "" {
		urlStr = u.Path
	}

	base := path.Base(urlStr)
	if base == "." || base == "/" || base == "" {
		return nil, fmt.Errorf("invalid url: %s", urlStr)
	}

	fullPath := filepath.Join(l.basePath, base)

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return f, nil
}
