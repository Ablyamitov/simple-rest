package mapper

import (
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/Ablyamitov/simple-rest/internal/store/web/dto"
)

func MapUserToDTO(user *entity.User) *dto.UserDTO {
	booksDTO := make([]*dto.BookDTO, len(user.Books))
	for i, book := range user.Books {
		booksDTO[i] = MapBookToDTO(book)
	}
	return &dto.UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Books:    booksDTO,
		Password: user.Password,
		Role:     user.Role,
	}
}

func MapDTOToUser(dto *dto.UserDTO) *entity.User {
	books := make([]*entity.Book, len(dto.Books))
	for i, bookDTO := range dto.Books {
		books[i] = MapDTOToBook(bookDTO)
	}

	return &entity.User{
		ID:       dto.ID,
		Name:     dto.Name,
		Email:    dto.Email,
		Books:    books,
		Password: dto.Password,
		Role:     dto.Role,
	}
}
