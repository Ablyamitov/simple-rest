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

type UserHandlerImpl struct {
	UserRepository repository.UserRepository
}

type UserHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	TakeBook(w http.ResponseWriter, r *http.Request)
	ReturnBook(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(userRepository repository.UserRepository) UserHandler {
	return &UserHandlerImpl{UserRepository: userRepository}
}

func (userHandler *UserHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := userHandler.UserRepository.GetAll(context.Background())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//wrapper.SendError(w, http.StatusNotFound, err, "UserHandlerImpl.GetAll")
			wrapper.LogError(err.Error(), "UserHandlerImpl.GetAll")
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.GetAll")
			wrapper.LogError(err.Error(), "UserHandlerImpl.GetAll")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	var usersDTO []*dto.UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, mapper.MapUserToDTO(&user))
	}
	if err := json.NewEncoder(w).Encode(usersDTO); err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.GetAll")
		wrapper.LogError(err.Error(), "UserHandlerImpl.GetAll")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (userHandler *UserHandlerImpl) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.GetByID")
		wrapper.LogError(err.Error(), "UserHandlerImpl.GetByID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := userHandler.UserRepository.GetByID(context.Background(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			//wrapper.SendError(w, http.StatusNotFound, err, "UserHandlerImpl.GetByID")
			wrapper.LogError(err.Error(), "UserHandlerImpl.GetByID")
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.GetByID")
			wrapper.LogError(err.Error(), "UserHandlerImpl.GetByID")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	userDTO := mapper.MapUserToDTO(user)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.GetByID")
		wrapper.LogError(err.Error(), "UserHandlerImpl.GetByID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (userHandler *UserHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	var userDTO *dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.Create")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.Validate(userDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.Create")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := mapper.MapDTOToUser(userDTO)
	err := userHandler.UserRepository.Create(context.Background(), user)
	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.Create")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	userDTO = mapper.MapUserToDTO(user)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.Create")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Create")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (userHandler *UserHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	var userDTO *dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.Update")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.Validate(userDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.Update")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := userHandler.UserRepository.Update(context.Background(), mapper.MapDTOToUser(userDTO))
	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.Update")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(mapper.MapUserToDTO(updatedUser)); err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.Update")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Update")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (userHandler *UserHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.Delete")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Delete")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = userHandler.UserRepository.Delete(context.Background(), id)
	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.Delete")
		wrapper.LogError(err.Error(), "UserHandlerImpl.Delete")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandlerImpl) TakeBook(w http.ResponseWriter, r *http.Request) {

	type TakeBookDTO struct {
		UserId int `json:"userId" validate:"required, notblank, gte=0"`
		BookId int `json:"bookId" validate:"required, notblank, gte=0"`
	}

	var takeBookDTO *TakeBookDTO

	if err := json.NewDecoder(r.Body).Decode(&takeBookDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.TakeBook")
		wrapper.LogError(err.Error(), "UserHandlerImpl.TakeBook")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.Validate(takeBookDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.TakeBook")
		wrapper.LogError(err.Error(), "UserHandlerImpl.TakeBook")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := userHandler.UserRepository.TakeBook(context.Background(), takeBookDTO.UserId, takeBookDTO.BookId)
	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.TakeBook")
		wrapper.LogError(err.Error(), "UserHandlerImpl.TakeBook")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (userHandler *UserHandlerImpl) ReturnBook(w http.ResponseWriter, r *http.Request) {

	type ReturnBookDTO struct {
		UserId int `json:"userId" validate:"required, notblank, gte=0"`
		BookId int `json:"bookId" validate:"required, notblank, gte=0"`
	}

	var returnBookDTO ReturnBookDTO

	if err := json.NewDecoder(r.Body).Decode(&returnBookDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.ReturnBook")
		wrapper.LogError(err.Error(), "UserHandlerImpl.ReturnBook")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.Validate(returnBookDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "UserHandlerImpl.ReturnBook")
		wrapper.LogError(err.Error(), "UserHandlerImpl.ReturnBook")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := userHandler.UserRepository.ReturnBook(context.Background(), returnBookDTO.UserId, returnBookDTO.BookId)
	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "UserHandlerImpl.ReturnBook")
		wrapper.LogError(err.Error(), "UserHandlerImpl.ReturnBook")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
