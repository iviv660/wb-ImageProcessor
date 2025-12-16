package image

type LocalStorage struct {
	basePath string
	baseURL  string
}

func New(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}
