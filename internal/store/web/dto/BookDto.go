package dto

type BookDTO struct {
	ID        int    `json:"id"`
	Title     string `json:"title" validate:"required,notblank"`
	Author    string `json:"author" validate:"required,notblank"`
	Available bool   `json:"available"`
}
