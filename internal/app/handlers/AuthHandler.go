package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Ablyamitov/simple-rest/internal/store/web/dto"
	"github.com/Ablyamitov/simple-rest/internal/store/web/mapper"
	"net/http"
	"time"

	"github.com/Ablyamitov/simple-rest/internal/app/utils"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"

	"github.com/dgrijalva/jwt-go"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	CheckAuth(w http.ResponseWriter, r *http.Request)
}

type AuthHandlerImpl struct {
	UserRepository repository.UserRepository
	Secret         string
}

type LoginRequest struct {
}

func NewAuthHandler(secret string, userRepository repository.UserRepository) AuthHandler {
	return &AuthHandlerImpl{UserRepository: userRepository, Secret: secret}
}

func (authHandler *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {

	var userDTO dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandlerImpl.Register")
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser entity.User
	existingUser, err := authHandler.UserRepository.GetByEmail(context.Background(), mapper.MapDTOToUser(&userDTO).Email)
	if err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandlerImpl.Register")
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if existingUser.ID != 0 {
		//wrapper.SendError(w, http.StatusBadRequest, errors.New("user with the same email already exists"), "AuthHandlerImpl.Register")
		wrapper.LogError(errors.New("user with the same email already exists").Error(), "AuthHandlerImpl.Register")
		http.Error(w, errors.New("user with the same email already exists").Error(), http.StatusBadRequest)
		return
	}

	userDTO.Password, err = utils.GenerateHashPassword(userDTO.Password)
	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandlerImpl.Register")
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = authHandler.UserRepository.Create(context.Background(), mapper.MapDTOToUser(&userDTO))
	if err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandlerImpl.Register")
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandlerImpl.Register")
	}
}

func (authHandler *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {

	var userDTO dto.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandlerImpl.Login")
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var existingUser entity.User
	existingUser, err := authHandler.UserRepository.GetByEmail(context.Background(), mapper.MapDTOToUser(&userDTO).Email)
	if err != nil {
		//wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandlerImpl.Login")
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if existingUser.ID == 0 {
		//wrapper.SendError(w, http.StatusBadRequest, errors.New("user with the same email is not exists"), "AuthHandlerImpl.Login")
		wrapper.LogError(errors.New("user with the same email is not exists").Error(), "AuthHandlerImpl.Login")
		http.Error(w, errors.New("user with the same email is not exists").Error(), http.StatusBadRequest)
		return
	}

	errHash := utils.CompareHashPassword(userDTO.Password, existingUser.Password)
	if !errHash {
		//wrapper.SendError(w, http.StatusBadRequest, errors.New("invalid password"), "AuthHandlerImpl.Login")
		wrapper.LogError(errors.New("invalid password").Error(), "AuthHandlerImpl.Login")
		http.Error(w, errors.New("invalid password").Error(), http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &entity.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(authHandler.Secret))

	if err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, errors.New("could not generate token"), "AuthHandlerImpl.Login")
		wrapper.LogError(errors.New("could not generate token").Error(), "AuthHandlerImpl.Login")
		http.Error(w, errors.New("could not generate token").Error(), http.StatusInternalServerError)
		return
	}

	bearerToken := "Bearer " + tokenString
	w.Header().Add("Authorization", bearerToken)
	w.Header().Add("role", claims.Role)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(existingUser); err != nil {
		//wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandlerImpl.Login")
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Login")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (authHandler *AuthHandlerImpl) CheckAuth(w http.ResponseWriter, r *http.Request) {

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		//wrapper.SendError(w, http.StatusUnauthorized, errors.New("authentication failed, because token is empty"), "AuthHandlerImpl.CheckAuth")
		wrapper.LogError(errors.New("authentication failed, because token is empty").Error(), "AuthHandlerImpl.CheckAuth")
		http.Error(w, errors.New("authentication failed, because token is empty").Error(), http.StatusUnauthorized)
		return
	}
	token := bearerToken[7:]

	claims, err := utils.ParseToken(token, authHandler.Secret)

	if err != nil {
		//wrapper.SendError(w, http.StatusUnauthorized, errors.New("token is not valid"), "AuthHandlerImpl.CheckAuth")
		wrapper.LogError(errors.New("token is not valid").Error(), "AuthHandlerImpl.CheckAuth")
		http.Error(w, errors.New("token is not valid").Error(), http.StatusUnauthorized)
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		//wrapper.SendError(w, http.StatusUnauthorized, errors.New("role does not have permission"), "AuthHandlerImpl.CheckAuth")
		wrapper.LogError(errors.New("role does not have permission").Error(), "AuthHandlerImpl.CheckAuth")
		http.Error(w, errors.New("role does not have permission").Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

}
