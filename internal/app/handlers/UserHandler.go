package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Ablyamitov/simple-rest/internal/store/web/dto"
	"github.com/Ablyamitov/simple-rest/internal/store/web/mapper"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"

	"github.com/Ablyamitov/simple-rest/internal/app/validation"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: userRepository}
}

func (userHandler *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := userHandler.UserRepository.GetAll(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.SendError(w, http.StatusNotFound, err, "UserHandler.GetAll")
		} else {
			wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.GetAll")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var usersDTO []*dto.UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, mapper.MapUserToDTO(&user))
	}
	if err := json.NewEncoder(w).Encode(usersDTO); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.GetAll")
	}
}

func (userHandler *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.GetByID")
		return
	}

	user, err := userHandler.UserRepository.GetByID(context.Background(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			wrapper.SendError(w, http.StatusNotFound, err, "UserHandler.GetByID")
		} else {
			wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.GetByID")
		}
		return
	}
	userDTO := mapper.MapUserToDTO(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.GetByID")
	}
}

func (userHandler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userDTO *dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.Create")
		return
	}

	if err := validation.Validate(userDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.Create")
		return
	}

	user := mapper.MapDTOToUser(userDTO)
	err := userHandler.UserRepository.Create(context.Background(), user)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.Create")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	userDTO = mapper.MapUserToDTO(user)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.Create")
	}
}

func (userHandler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var userDTO *dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.Update")
		return
	}

	if err := validation.Validate(userDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.Update")
		return
	}

	updatedUser, err := userHandler.UserRepository.Update(context.Background(), mapper.MapDTOToUser(userDTO))
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.Update")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(mapper.MapUserToDTO(updatedUser)); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.Update")
	}

}

func (userHandler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.Delete")
		return
	}
	err = userHandler.UserRepository.Delete(context.Background(), id)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.Delete")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) TakeBook(w http.ResponseWriter, r *http.Request) {

	type TakeBookDTO struct {
		UserId int `json:"userId" validate:"required, notblank, gte=0"`
		BookId int `json:"bookId" validate:"required, notblank, gte=0"`
	}

	var takeBookDTO *TakeBookDTO

	if err := json.NewDecoder(r.Body).Decode(&takeBookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.TakeBook")
		return
	}

	if err := validation.Validate(takeBookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.TakeBook")
		return
	}

	err := userHandler.UserRepository.TakeBook(context.Background(), takeBookDTO.UserId, takeBookDTO.BookId)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.TakeBook")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandler) ReturnBook(w http.ResponseWriter, r *http.Request) {

	type ReturnBookDTO struct {
		UserId int `json:"userId" validate:"required, notblank, gte=0"`
		BookId int `json:"bookId" validate:"required, notblank, gte=0"`
	}

	var returnBookDTO ReturnBookDTO

	if err := json.NewDecoder(r.Body).Decode(&returnBookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.ReturnBook")
		return
	}

	if err := validation.Validate(returnBookDTO); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "UserHandler.ReturnBook")
		return
	}

	err := userHandler.UserRepository.ReturnBook(context.Background(), returnBookDTO.UserId, returnBookDTO.BookId)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandler.ReturnBook")
		return
	}
	w.WriteHeader(http.StatusOK)
}
