package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/Ablyamitov/simple-rest/internal/store/web/dto"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ablyamitov/simple-rest/internal/store/db/repository/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBookHandler_GetAll(t *testing.T) {

	type mockBehavior func(mockRepository *repository.MockBookRepository)
	testCases := []struct {
		name               string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBooks      []dto.BookDTO
	}{
		{
			name: "Test 1: OK",
			mockBehavior: func(mockRepository *repository.MockBookRepository) {
				templateBooks := []entity.Book{
					{
						ID:        1,
						Title:     "Test english",
						Author:    "Test author",
						Available: true,
					},
					{
						ID:        2,
						Title:     "Test spanish",
						Author:    "Test author",
						Available: true,
					},
				}
				mockRepository.EXPECT().GetALL(gomock.Any()).Return(templateBooks, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBooks: []dto.BookDTO{
				{
					ID:        1,
					Title:     "Test english",
					Author:    "Test author",
					Available: true,
				},
				{
					ID:        2,
					Title:     "Test spanish",
					Author:    "Test author",
					Available: true,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := repository.NewMockBookRepository(ctrl)
			testCase.mockBehavior(mockRepository)
			handler := NewBookHandler(mockRepository)

			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			w := httptest.NewRecorder()
			handler.GetAll(w, req)

			resp := w.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatalf("could not close resp result: %v", err)
				}
			}(resp.Body)
			var responseBooks []dto.BookDTO
			err := json.NewDecoder(resp.Body).Decode(&responseBooks)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, testCase.expectedBooks, responseBooks)
		})
	}
}

func TestBookHandler_GetById(t *testing.T) {

	type mockBehavior func(mockRepository *repository.MockBookRepository)

	testCases := []struct {
		name               string
		inputID            string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       dto.BookDTO
		responseBody       dto.BookDTO
	}{
		{
			name:    "Test 1: OK",
			inputID: "1",
			mockBehavior: func(mockRepository *repository.MockBookRepository) {
				templateBook := entity.Book{
					ID:        1,
					Title:     "Test english",
					Author:    "Test author",
					Available: true,
				}
				mockRepository.EXPECT().GetByID(gomock.Any(), gomock.Eq(1)).Return(&templateBook, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: dto.BookDTO{
				ID:        1,
				Title:     "Test english",
				Author:    "Test author",
				Available: true,
			},
		},
		{
			name:               "Test 2: Invalid ID",
			inputID:            "abc",
			mockBehavior:       func(mockRepository *repository.MockBookRepository) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Test 3: Not Found",
			inputID: "999",
			mockBehavior: func(mockRepository *repository.MockBookRepository) {
				mockRepository.EXPECT().GetByID(gomock.Any(), gomock.Eq(999)).Return(nil, errors.New("book not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       dto.BookDTO{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepository := repository.NewMockBookRepository(ctrl)
			testCase.mockBehavior(mockRepository)
			handler := NewBookHandler(mockRepository)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%s", testCase.inputID), nil)
			req, err := func(req *http.Request, strId string) (*http.Request, error) {
				chiCtx := chi.NewRouteContext()
				chiCtx.URLParams.Add("id", strId)
				return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)), nil
			}(req, testCase.inputID)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			w := httptest.NewRecorder()
			handler.GetById(w, req)

			resp := w.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatalf("could not close resp result: %v", err)
				}
			}(resp.Body)

			err = json.NewDecoder(resp.Body).Decode(&testCase.responseBody)

			if testCase.expectedStatusCode == http.StatusOK {
				err := json.NewDecoder(resp.Body).Decode(&testCase.responseBody)
				if err != io.EOF {
					assert.NoError(t, err)
				}
				assert.Equal(t, testCase.expectedBody, testCase.responseBody)
			} else {
				assert.Equal(t, testCase.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestBookHandler_Create(t *testing.T) {

}
