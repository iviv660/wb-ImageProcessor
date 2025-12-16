package v1

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}

	img, err := a.service.GetImage(r.Context(), id)
	if err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	if img.ChangedImageURL == nil || *img.ChangedImageURL == "" {
		http.Error(w, "image is not processed yet", http.StatusConflict) // 409
		return
	}

	reader, err := a.storage.Open(r.Context(), *img.ChangedImageURL)
	if err != nil {
		http.Error(w, "failed open image", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = io.Copy(w, reader)
}
