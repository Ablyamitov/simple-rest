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

type BookHandler struct {
	BookRepository repository.BookRepository
}

func NewBookHandler(bookRepository repository.BookRepository) *BookHandler {
	return &BookHandler{BookRepository: bookRepository}
}

func (bookHandler *BookHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	books, err := bookHandler.BookRepository.GetALL(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.SendError(w, http.StatusNotFound, err, "BookHandler.GetAll")
		} else {
			wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.GetAll")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var booksDTO []*dto.BookDTO
	for _, book := range books {
		booksDTO = append(booksDTO, mapper.MapBookToDTO(&book))
	}
	if err := json.NewEncoder(w).Encode(booksDTO); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.GetAll")
		return
	}
}

func (bookHandler *BookHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		wrapper.SendError(w, http.StatusBadRequest, err, "BookHandler.GetByID")
		return
	}
	book, err := bookHandler.BookRepository.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.SendError(w, http.StatusNotFound, err, "BookHandler.GetById")
		} else {
			wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.GetById")
		}
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mapper.MapBookToDTO(book)); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.GetById")
		return
	}

}

func (bookHandler *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var bookDTO *dto.BookDTO
	if err := json.NewDecoder(r.Body).Decode(&bookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "BookHandler.Create")
		return
	}

	if err := validation.Validate(bookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "BookHandler.Create")
		return
	}

	err := bookHandler.BookRepository.Create(context.Background(), mapper.MapDTOToBook(bookDTO))
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.Create")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(bookDTO); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.Create")
		return
	}
}

func (bookHandler *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var bookDTO *dto.BookDTO
	if err := json.NewDecoder(r.Body).Decode(&bookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "BookHandler.Update")
		return
	}

	if err := validation.Validate(bookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "BookHandler.Update")
		return
	}

	updatedBook, err := bookHandler.BookRepository.Update(context.Background(), mapper.MapDTOToBook(bookDTO))
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.Update")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mapper.MapBookToDTO(updatedBook)); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.Update")
		return
	}

}

func (bookHandler *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "BookHandler.Delete")
		return
	}

	err = bookHandler.BookRepository.Delete(context.Background(), id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "BookHandler.Delete")
		return
	}

}
