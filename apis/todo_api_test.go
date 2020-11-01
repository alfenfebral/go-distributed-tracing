package apis_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"../apis"
	mockServices "../mocks/services"
	"../models"
	"../utils"

	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/faker"
)

// TestGetAllSuccess - testing GetAll [200]
func TestGetAllSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockService := new(mockServices.TodoService)
	mockListTodo := make([]*models.Todo, 0)
	mockListTodo = append(mockListTodo, &mockTodo)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "/todo", nil)
	assert.NoError(t, err)

	q := req.URL.Query()
	q.Add("q", "")
	q.Add("page", "1")
	q.Add("per_page", "10")

	// Encode the query to support req.URL.String()
	req.URL.RawQuery = q.Encode()

	currentPage := utils.CurrentPage(q.Get("page"))
	perPage := utils.PerPage(q.Get("per_page"))
	offset := utils.Offset(currentPage, perPage)

	// Mock GetAll service
	mockService.On("GetAll", q.Get("q"), perPage, offset).Return(mockListTodo, len(mockListTodo), nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.GetAll)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the mock called
	mockService.AssertExpectations(t)
}

// TestGetAllFailedValidationError - testing GetAll [400]
// failed validation
func TestGetAllFailedValidationError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockService := new(mockServices.TodoService)
	mockListTodo := make([]*models.Todo, 0)
	mockListTodo = append(mockListTodo, &mockTodo)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "/todo", nil)
	assert.NoError(t, err)

	q := req.URL.Query()
	q.Add("q", "")
	q.Add("page", "-1")
	q.Add("per_page", "-10")

	// Encode the query to support req.URL.String()
	req.URL.RawQuery = q.Encode()

	// Mock GetAll service
	mockService.On("GetAll", "", 10, utils.Offset(1, 10)).Return(mockListTodo, len(mockListTodo), nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.GetAll)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestGetAllFailedDatabaseError - testing GetAll [500]
func TestGetAllFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockService := new(mockServices.TodoService)
	mockListTodo := make([]*models.Todo, 0)
	mockListTodo = append(mockListTodo, &mockTodo)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "/todo", nil)
	assert.NoError(t, err)

	q := req.URL.Query()
	q.Add("q", "")
	q.Add("page", "1")
	q.Add("per_page", "10")

	// Encode the query to support req.URL.String()
	req.URL.RawQuery = q.Encode()

	currentPage := utils.CurrentPage(q.Get("page"))
	perPage := utils.PerPage(q.Get("per_page"))
	offset := utils.Offset(currentPage, perPage)

	// Mock GetAll service
	mockService.On("GetAll", q.Get("q"), perPage, offset).Return(mockListTodo, len(mockListTodo), errors.New(""))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.GetAll)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check if the mock called
	mockService.AssertExpectations(t)
}

// TestGetByIdSuccess - testing GetById [200]
func TestGetByIdSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	// Todo id
	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), nil)
	assert.NoError(t, err)

	// Mock GetByID service
	mockService.On("GetByID", "").Return(&mockTodo, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.GetByID)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

func TestGetByIdNotFoundError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	// Todo id
	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), nil)
	assert.NoError(t, err)

	// Mock GetByID service
	mockService.On("GetByID", "").Return(&mockTodo, errors.New("not found"))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.GetByID)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

func TestGetByIdDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	// Todo id
	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), nil)
	assert.NoError(t, err)

	// Mock GetByID service
	mockService.On("GetByID", "").Return(&mockTodo, errors.New(""))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.GetByID)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}
