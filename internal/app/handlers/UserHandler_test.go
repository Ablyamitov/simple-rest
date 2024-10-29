package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockUserRepository(ctrl)

	handler := NewUserHandler(mockRepository)

	users := []entity.User{
		{ID: 1, Name: "John", Email: "john@example.com"},
		{ID: 2, Name: "Alex", Email: "alex@example.com"},
	}

	mockRepository.EXPECT().
		GetAll(gomock.Any()).
		Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseUsers []entity.User
	err := json.NewDecoder(resp.Body).Decode(&responseUsers)
	assert.NoError(t, err)

	assert.Equal(t, users, responseUsers)

	mockRepository.EXPECT().
		GetAll(gomock.Any()).
		Return(nil, errors.New("internal server error"))

	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	w = httptest.NewRecorder()

	handler.GetAll(w, req)

	resp = w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUserHandler_GetById(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockUserRepository(ctrl)

	handler := NewUserHandler(mockRepository)

	//1
	user := &entity.User{
		ID:       1,
		Name:     "John",
		Email:    "john@example.com",
		Password: "1234",
		Role:     "user",
	}

	mockRepository.EXPECT().
		GetByID(gomock.Any(), gomock.Eq(1)).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = chiCtxWithID(req, 1)

	w := httptest.NewRecorder()

	handler.GetById(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseUser entity.User
	err := json.NewDecoder(resp.Body).Decode(&responseUser)
	assert.NoError(t, err)

	assert.Equal(t, *user, responseUser)

	//2
	req = httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	w = httptest.NewRecorder()

	handler.GetById(w, req)

	resp = w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	//3
	mockRepository.EXPECT().
		GetByID(gomock.Any(), gomock.Eq(2)).
		Return(nil, errors.New("internal server error"))

	req = httptest.NewRequest(http.MethodGet, "/users/2", nil)
	req = chiCtxWithID(req, 2)
	w = httptest.NewRecorder()

	handler.GetById(w, req)

	resp = w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestUserHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockUserRepository(ctrl)

	handler := NewUserHandler(mockRepository)

	newUser := &entity.User{
		ID:       1,
		Name:     "John",
		Email:    "John@example.com",
		Password: "1234",
		Role:     "user",
	}
	userJSON, err := json.Marshal(newUser)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/add", bytes.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepository.EXPECT().
		Create(gomock.Any(), gomock.Eq(newUser)).
		Return(nil)

	handler.Create(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	invalidUser := entity.User{
		ID:       2,
		Name:     "",
		Email:    "Alex@example.com",
		Password: "1234",
		Role:     "user",
	}

	invalidUserJSON, err := json.Marshal(invalidUser)
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(invalidUserJSON))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()

	handler.Create(w, req)

	resp = w.Result()

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUserHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockUserRepository(ctrl)
	handler := NewUserHandler(mockRepository)

	existingUser := &entity.User{
		ID:       1,
		Name:     "John",
		Email:    "john@example.com",
		Password: "1234",
		Role:     "user",
	}

	updatedUser := &entity.User{
		ID:       1,
		Name:     "Alex",
		Email:    "alex@example.com",
		Password: "1234",
		Role:     "user",
	}

	userJSON, err := json.Marshal(updatedUser)
	assert.NoError(t, err)
	assert.NotEqual(t, *existingUser, *updatedUser)

	req := httptest.NewRequest(http.MethodPatch, "/users/update", bytes.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepository.EXPECT().
		Update(gomock.Any(), gomock.Eq(updatedUser)).
		Return(updatedUser, nil)

	handler.Update(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)

	var responseUser entity.User
	err = json.NewDecoder(resp.Body).Decode(&responseUser)
	assert.NoError(t, err)

	assert.Equal(t, *updatedUser, responseUser)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestUserHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockUserRepository(ctrl)
	handler := NewUserHandler(mockRepository)
	//1
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req = chiCtxWithID(req, 1)

	w := httptest.NewRecorder()

	mockRepository.EXPECT().
		Delete(gomock.Any(), gomock.Eq(1)).
		Return(nil)

	handler.Delete(w, req)
	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	//2
	req = httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
	w = httptest.NewRecorder()

	handler.Delete(w, req)

	resp = w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	//3
	mockRepository.EXPECT().
		Delete(gomock.Any(), gomock.Eq(2)).
		Return(sql.ErrNoRows)

	req = httptest.NewRequest(http.MethodDelete, "/users/2", nil)
	req = chiCtxWithID(req, 2) // Добавляем ID в контекст запроса
	w = httptest.NewRecorder()

	handler.Delete(w, req)

	resp = w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatalf("could not close resp result: %v", err)
		}
	}(resp.Body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func chiCtxWithID(req *http.Request, id int) *http.Request {
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", strconv.Itoa(id))
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
}
