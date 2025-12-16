package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/iviv660/wb-ImageProcessor/internal/service"
	"github.com/iviv660/wb-ImageProcessor/internal/storage"
)

type API struct {
	mux     *chi.Mux
	service service.ImageService
	storage storage.Storage
}

func New(mux *chi.Mux, imageService service.ImageService, storage storage.Storage) *API {
	return &API{
		mux:     mux,
		service: imageService,
		storage: storage,
	}
}

func (api *API) Register() {
	api.mux.Post("/upload", api.Upload)
	api.mux.Get("/image/{id}", api.Get)
	api.mux.Delete("/image/{id}", api.Delete)
}
