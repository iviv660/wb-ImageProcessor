package model

type Image struct {
	ID               int64   `db:"id"`
	FileName         string  `db:"file_name"`
	Format           string  `db:"format"`
	OriginalImageURL string  `db:"original_url"`
	ChangedImageURL  *string `db:"changed_url"`
}
