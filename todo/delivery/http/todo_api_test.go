package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "go-distributed-tracing/todo/delivery/http"
	"go-distributed-tracing/todo/models"
	"go-distributed-tracing/utils"

	mockServices "go-distributed-tracing/todo/mocks/services"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/sdk/trace"
)

var ErrDefault error = errors.New("error")
var ErrNotFound error = errors.New("not found")
var WhenError400EOF string = "when return 400 bad request (error EOF)"
var WhenError500Service string = "when return 500 internal error (error service)"
var WhenError500Query string = "when return 500 internal error (error query)"
var WhenError400Validation string = "when return 400 bad request (error validation)"
var WhenError404NotFound string = "when return 404 not found (resouce not found)"
var WhenSuccess201Created string = "when return 201 created"
var WhenSuccess200OK string = "when return 200 ok"

func TestNewTodoHTTPHandler(t *testing.T) {
	utils.InitializeValidator()

	router := chi.NewRouter()
	mockService := new(mockServices.TodoService)

	mockService.On(
		"Create",
		mock.AnythingOfType("context.Context"),
		mock.AnythingOfType("*models.Todo"),
	).Return(&models.Todo{}, nil)

	handlers.NewTodoHTTPHandler(router, trace.NewTracerProvider(), mockService)
}

// TestTodoGetAll - testing GetAll [200]
func TestTodoGetAll(t *testing.T) {
	t.Run(WhenError400Validation, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=-1&per_page=-1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=1&per_page=10", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On(
			"GetAll",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(nil, 1, ErrDefault)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		mockListTodo := make([]*models.Todo, 0)
		mockListTodo = append(mockListTodo, &models.Todo{})

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=1&per_page=10", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On(
			"GetAll",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(mockListTodo, 1, nil)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestTodoCreate - testing create [201]
func TestTodoCreate(t *testing.T) {
	t.Run(WhenError400EOF, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader([]byte("")))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("when return 400 bad request (error validation) ", func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "",
			"description": "",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("when error 500 internal error (error service)", func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "lorem ipsum",
			"description": "desc",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil, errors.New(""))
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run("when return 201 created", func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "lorem ipsum",
			"description": "desc",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestTodoGetByID - testing GetByID [200]
func TestTodoGetByID(t *testing.T) {
	t.Run(WhenError404NotFound, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, ErrNotFound)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, ErrDefault)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&models.Todo{}, nil)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestTodoUpdate - testing update [200]
func TestTodoUpdate(t *testing.T) {
	t.Run(WhenError400EOF, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodPut, "/api/v1/product?id=1", bytes.NewReader([]byte("")))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(&models.Todo{}, nil)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run(WhenError400Validation, func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "",
			"description": "",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError404NotFound, func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(nil, ErrNotFound)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(nil, ErrDefault)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(&models.Todo{}, nil)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}

// TestDeleteSuccess - testing delete [200]
func TestTodoDelete(t *testing.T) {
	t.Run(WhenError404NotFound, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/product?id=", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(ErrNotFound)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(ErrDefault)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		utils.InitializeValidator()

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockService := new(mockServices.TodoService)
		mockService.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)
		tp := trace.NewTracerProvider()

		todoHandler := handlers.NewTodoHTTPHandler(chi.NewRouter(), tp, mockService)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockService.AssertExpectations(t)
	})
}
