package mapper

import (
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/Ablyamitov/simple-rest/internal/store/web/dto"
)

func MapBookToDTO(book *entity.Book) *dto.BookDTO {
	return &dto.BookDTO{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Available: book.Available,
	}
}

func MapDTOToBook(dto *dto.BookDTO) *entity.Book {
	return &entity.Book{
		ID:        dto.ID,
		Title:     dto.Title,
		Author:    dto.Author,
		Available: dto.Available,
	}
}
