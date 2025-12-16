package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/iviv660/wb-ImageProcessor/internal/service"
)

func (a *API) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "bad multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "bad form file", http.StatusBadRequest)
		return
	}

	format := strings.TrimPrefix(strings.ToLower(header.Filename), ".")
	image, err := a.service.Upload(r.Context(), service.UploadInput{
		FileName: header.Filename,
		Format:   format,
		Size:     header.Size,
		Data:     file,
	})
	if err != nil {
		http.Error(w, "failed upload", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id": image.ID,
	})
}
