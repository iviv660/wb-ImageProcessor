package v1

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}

	if err := a.service.DeleteImage(r.Context(), id); err != nil {
		http.Error(w, "Failed delete image", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
