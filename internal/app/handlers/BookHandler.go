package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Ablyamitov/simple-rest/internal/app/validation"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"
	"github.com/Ablyamitov/simple-rest/internal/store/web/dto"
	"github.com/Ablyamitov/simple-rest/internal/store/web/mapper"
	"github.com/jackc/pgx/v5"

	"github.com/go-chi/chi/v5"
)

type BookHandlerImpl struct {
	BookRepository repository.BookRepository
}

type BookHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewBookHandler(bookRepository repository.BookRepository) BookHandler {
	return &BookHandlerImpl{BookRepository: bookRepository}
}

func (bookHandler *BookHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	books, err := bookHandler.BookRepository.GetALL(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.LogError(err.Error(), "BookHandlerImpl.GetAll")
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			wrapper.LogError(err.Error(), "BookHandlerImpl.GetAll")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	var booksDTO []*dto.BookDTO
	for _, book := range books {
		booksDTO = append(booksDTO, mapper.MapBookToDTO(&book))
	}
	if err := json.NewEncoder(w).Encode(booksDTO); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.GetAll")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (bookHandler *BookHandlerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.GetByID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	book, err := bookHandler.BookRepository.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.LogError(err.Error(), "BookHandlerImpl.GetByID")
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			wrapper.LogError(err.Error(), "BookHandlerImpl.GetByID")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mapper.MapBookToDTO(book)); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.GetByID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (bookHandler *BookHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	var bookDTO *dto.BookDTO
	if err := json.NewDecoder(r.Body).Decode(&bookDTO); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.Validate(bookDTO); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := bookHandler.BookRepository.Create(context.Background(), mapper.MapDTOToBook(bookDTO))
	if err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(bookDTO); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (bookHandler *BookHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	var bookDTO *dto.BookDTO
	if err := json.NewDecoder(r.Body).Decode(&bookDTO); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.Validate(bookDTO); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedBook, err := bookHandler.BookRepository.Update(context.Background(), mapper.MapDTOToBook(bookDTO))
	if err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mapper.MapBookToDTO(updatedBook)); err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (bookHandler *BookHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Delete")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bookHandler.BookRepository.Delete(context.Background(), id)

	w.WriteHeader(http.StatusOK)
	if err != nil {
		wrapper.LogError(err.Error(), "BookHandlerImpl.Delete")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
