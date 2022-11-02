package apis_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-clean-architecture/apis"
	mockServices "go-clean-architecture/mocks/services"
	"go-clean-architecture/models"
	"go-clean-architecture/utils"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bxcodec/faker"
)

func TestNewTodoHTTPHandler(t *testing.T) {
	router := chi.NewRouter()
	mockService := new(mockServices.TodoService)

	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)
	mockListTodo := make([]*models.Todo, 0)
	mockListTodo = append(mockListTodo, &mockTodo)

	mockService.On("GetAll", "", 10, utils.Offset(1, 10)).Return(mockListTodo, len(mockListTodo), nil)
	mockService.On("GetByID", "").Return(&mockTodo, nil)
	mockService.On("Create", mock.Anything).Return(&mockTodo, nil)
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil, nil)
	mockService.On("Delete", mock.Anything).Return(nil)

	apis.NewTodoHTTPHandler(router, mockService)
}

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

func TestGetByIdFailedDatabaseError(t *testing.T) {
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

// TestCreateSuccess - testing create [200]
func TestCreateSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "Lorem Ipsum",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "todo", bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Create", mock.Anything).Return(&mockTodo, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Create)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestCreateFailedValidationError - testing create [400]
func TestCreateFailedValidationError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "todo", bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Create", mock.Anything).Return(&mockTodo, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Create)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestCreateFailedEmptyBody - testing create [400]
func TestCreateFailedEmptyBody(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "todo", strings.NewReader(""))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Create", mock.Anything).Return(&mockTodo, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Create)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestCreateFailedDatabaseError - testing create [500]
func TestCreateFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "Lorem Ipsum",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "todo", bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Create", mock.Anything).Return(&mockTodo, errors.New(""))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Create)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestUpdateSuccess - testing update [200]
func TestUpdateSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "Lorem Ipsum",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Update)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestUpdateFailedValidationError - testing update [400]
func TestUpdateFailedValidationError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Update)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestUpdateFailedEmptyBody - testing update [400]
func TestUpdateFailedEmptyBody(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), strings.NewReader(""))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil, nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Update)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestUpdateFailedNotFound - testing update [404]
func TestUpdateFailedNotFound(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "Lorem Ipsum",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Update)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestUpdateFailedDatabaseError - testing update [404]
func TestUpdateFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	// Mock TodoService
	mockService := new(mockServices.TodoService)

	mockPostBody := map[string]interface{}{
		"title":       "Lorem Ipsum",
		"description": "Lorem Ipsum",
	}
	body, _ := json.Marshal(mockPostBody)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock Create service
	// mocking dynamic argument to exact with return value
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Update)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestDeleteSuccess - testing create [200]
func TestDeleteSuccess(t *testing.T) {
	// Mock TodoService
	mockService := new(mockServices.TodoService)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock service function
	mockService.On("Delete", mock.Anything).Return(nil)

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Delete)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestDeleteFailedNotFound - testing create [404]
func TestDeleteFailedNotFound(t *testing.T) {
	// Mock TodoService
	mockService := new(mockServices.TodoService)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock service function
	mockService.On("Delete", mock.Anything).Return(errors.New("not found"))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Delete)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}

// TestDeleteFailedDatabaseError - testing create [404]
func TestDeleteFailedDatabaseError(t *testing.T) {
	// Mock TodoService
	mockService := new(mockServices.TodoService)

	id := "5f9e9dc1c5e4cd34e6ac0bb7"

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("todo/%s", id), nil)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Mock service function
	mockService.On("Delete", mock.Anything).Return(errors.New(""))

	// Mock service in handler
	todoHandler := apis.Todohandler{
		TodoService: mockService,
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.Delete)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what expected
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check if the mock called
	// mockService.AssertExpectations(t)
}
