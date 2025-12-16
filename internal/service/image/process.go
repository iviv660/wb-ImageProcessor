package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"

	"github.com/iviv660/wb-ImageProcessor/internal/model"
)

func processBytes(b []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)

	wmW, wmH := 200, 60
	padding := 20
	x0 := rgba.Bounds().Max.X - wmW - padding
	y0 := rgba.Bounds().Max.Y - wmH - padding
	rect := image.Rect(x0, y0, x0+wmW, y0+wmH)

	draw.Draw(rgba, rect, &image.Uniform{C: color.RGBA{0, 0, 0, 120}}, image.Point{}, draw.Over)

	out := new(bytes.Buffer)
	switch format {
	case "png":
		_ = png.Encode(out, rgba)
	default:
		_ = jpeg.Encode(out, rgba, &jpeg.Options{Quality: 85})
	}
	return out.Bytes(), nil
}

func (s *Service) ProcessImage(ctx context.Context, id int64) (any, error) {
	if id <= 0 {
		return model.Image{}, fmt.Errorf("invalid id")
	}

	img, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Image{}, fmt.Errorf("get image: %w", err)
	}
	if img.OriginalImageURL == "" {
		return model.Image{}, fmt.Errorf("original url is empty")
	}

	rc, err := s.storage.Open(ctx, img.OriginalImageURL)
	if err != nil {
		return model.Image{}, fmt.Errorf("open original: %w", err)
	}
	defer rc.Close()

	origBytes, err := io.ReadAll(rc)
	if err != nil {
		return model.Image{}, fmt.Errorf("read original: %w", err)
	}

	changedBytes, err := processBytes(origBytes)
	if err != nil {
		return model.Image{}, fmt.Errorf("process image: %w", err)
	}

	safeName := filepath.Base(img.FileName)
	changedName := fmt.Sprintf("changed_%d_%s", img.ID, safeName)

	changedURL, err := s.storage.Save(ctx, changedName, bytes.NewReader(changedBytes))
	if err != nil {
		return model.Image{}, fmt.Errorf("save changed: %w", err)
	}

	if img.ChangedImageURL != nil && *img.ChangedImageURL != changedURL {
		_ = s.storage.Delete(ctx, *img.ChangedImageURL)
	}

	img.ChangedImageURL = &changedURL
	updated, err := s.repo.Update(ctx, img)
	if err != nil {
		_ = s.storage.Delete(ctx, changedURL)
		return model.Image{}, fmt.Errorf("update db: %w", err)
	}

	return updated, nil
}
