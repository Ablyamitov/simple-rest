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

var (
	errNotUniqueEmail        = errors.New("user with the same email already exists")
	errNotFoundWithSameEmail = errors.New("user with the same email is not exists")
	errInvalidPassword       = errors.New("invalid password")
	errGenerateToken         = errors.New("could not generate token")
	errEmptyToken            = errors.New("authentication failed, because token is empty")
	errNotValidToken         = errors.New("token is not valid")
	errAccessDenied          = errors.New("role does not have permission")
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
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser entity.User
	existingUser, err := authHandler.UserRepository.GetByEmail(context.Background(), mapper.MapDTOToUser(&userDTO).Email)
	if err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if existingUser.ID != 0 {
		wrapper.LogError(errNotUniqueEmail.Error(), "AuthHandlerImpl.Register")
		http.Error(w, errNotUniqueEmail.Error(), http.StatusBadRequest)
		return
	}

	userDTO.Password, err = utils.GenerateHashPassword(userDTO.Password)
	if err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = authHandler.UserRepository.Create(context.Background(), mapper.MapDTOToUser(&userDTO))
	if err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Register")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (authHandler *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {

	var userDTO dto.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var existingUser entity.User
	existingUser, err := authHandler.UserRepository.GetByEmail(context.Background(), mapper.MapDTOToUser(&userDTO).Email)
	if err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Login")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if existingUser.ID == 0 {
		wrapper.LogError(errNotFoundWithSameEmail.Error(), "AuthHandlerImpl.Login")
		http.Error(w, errNotFoundWithSameEmail.Error(), http.StatusBadRequest)
		return
	}

	errHash := utils.CompareHashPassword(userDTO.Password, existingUser.Password)
	if !errHash {
		wrapper.LogError(errInvalidPassword.Error(), "AuthHandlerImpl.Login")
		http.Error(w, errInvalidPassword.Error(), http.StatusBadRequest)
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
		wrapper.LogError(errGenerateToken.Error(), "AuthHandlerImpl.Login")
		http.Error(w, errGenerateToken.Error(), http.StatusInternalServerError)
		return
	}

	bearerToken := "Bearer " + tokenString
	w.Header().Add("Authorization", bearerToken)
	w.Header().Add("role", claims.Role)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(existingUser); err != nil {
		wrapper.LogError(err.Error(), "AuthHandlerImpl.Login")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (authHandler *AuthHandlerImpl) CheckAuth(w http.ResponseWriter, r *http.Request) {

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		wrapper.LogError(errEmptyToken.Error(), "AuthHandlerImpl.CheckAuth")
		http.Error(w, errEmptyToken.Error(), http.StatusUnauthorized)
		return
	}
	token := bearerToken[7:]

	claims, err := utils.ParseToken(token, authHandler.Secret)

	if err != nil {
		wrapper.LogError(errNotValidToken.Error(), "AuthHandlerImpl.CheckAuth")
		http.Error(w, errNotValidToken.Error(), http.StatusUnauthorized)
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		wrapper.LogError(errAccessDenied.Error(), "AuthHandlerImpl.CheckAuth")
		http.Error(w, errAccessDenied.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

}
