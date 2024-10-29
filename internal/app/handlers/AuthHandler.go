package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Ablyamitov/simple-rest/internal/app/utils"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"

	"github.com/dgrijalva/jwt-go"
)

type AuthHandler struct {
	UserRepository repository.UserRepository
	Secret         string
}
type RegisterRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,notblank"`
	Email    string `json:"email" validate:"email,required,notblank"`
	Password string `json:"password"`
}

type LoginRequest struct {
}

func NewAuthHandler(secret string, userRepository repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepository: userRepository, Secret: secret}
}

func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandler.Register")
		return
	}

	var existingUser entity.User
	existingUser, err := authHandler.UserRepository.GetByEmail(context.Background(), user.Email)
	if err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandler.Register")
		return
	}
	if existingUser.ID != 0 {
		wrapper.SendError(w, http.StatusBadRequest, errors.New("user with the same email already exists"), "AuthHandler.Register")
		return
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandler.Register")
		return
	}
	err = authHandler.UserRepository.Create(context.Background(), &user)
	if err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandler.Register")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandler.Register")
	}
}

func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var user entity.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandler.Login")
		return
	}
	var existingUser entity.User
	existingUser, err := authHandler.UserRepository.GetByEmail(context.Background(), user.Email)
	if err != nil {
		wrapper.SendError(w, http.StatusBadRequest, err, "AuthHandler.Login")
		return
	}
	if existingUser.ID == 0 {
		wrapper.SendError(w, http.StatusBadRequest, errors.New("user with the same email is not exists"), "AuthHandler.Login")
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)
	if !errHash {
		wrapper.SendError(w, http.StatusBadRequest, errors.New("invalid password"), "AuthHandler.Login")
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
		wrapper.SendError(w, http.StatusInternalServerError, errors.New("could not generate token"), "AuthHandler.Login")
		return
	}

	bearerToken := "Bearer " + tokenString
	w.Header().Add("Authorization", bearerToken)
	w.Header().Add("role", claims.Role)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(existingUser); err != nil {
		wrapper.SendError(w, http.StatusInternalServerError, err, "AuthHandler.Login")
	}
}

func (authHandler *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		wrapper.SendError(w, http.StatusUnauthorized, errors.New("authentication failed, because token is empty"), "AuthHandler.CheckAuth")
		return
	}
	token := bearerToken[7:]

	claims, err := utils.ParseToken(token, authHandler.Secret)

	if err != nil {
		wrapper.SendError(w, http.StatusUnauthorized, errors.New("token is not valid"), "AuthHandler.CheckAuth")
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		wrapper.SendError(w, http.StatusUnauthorized, errors.New("role does not have permission"), "AuthHandler.CheckAuth")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

}
