package services_test

import (
	"context"
	"errors"
	mockRepositories "go-distributed-tracing/todo/mocks/repository"
	"go-distributed-tracing/todo/models"
	"go-distributed-tracing/todo/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ErrDefault error = errors.New("error")
var DefaultID string = "1"

func TestTodoGetAll(t *testing.T) {
	t.Run("success when find all", func(t *testing.T) {
		mockList := make([]*models.Todo, 0)
		mockList = append(mockList, &models.Todo{})

		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On(
			"FindAll",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"), mock.AnythingOfType("int"),
		).Return(mockList, nil)
		mockRepository.On(
			"CountFindAll",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(10, nil)

		ctx := context.Background()
		results, count, err := service.GetAll(ctx, "keyword", 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, count, 10)
		assert.Equal(t, mockList, results)
	})

	t.Run("error when find all", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On(
			"FindAll",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(nil, ErrDefault)
		mockRepository.On(
			"CountFindAll",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(10, nil)

		ctx := context.Background()
		results, count, err := service.GetAll(ctx, "keyword", 10, 0)

		assert.Nil(t, results)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
	})

	t.Run("error when count find all", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On(
			"FindAll",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),
		).Return(nil, nil)
		mockRepository.On(
			"CountFindAll",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(10, ErrDefault)

		ctx := context.Background()
		results, count, err := service.GetAll(ctx, "keyword", 10, 0)

		assert.Nil(t, results)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
	})
}

func TestTodoGetByID(t *testing.T) {
	t.Run("success when find by id", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On("FindById", mock.Anything, mock.AnythingOfType("string")).Return(mockTodo, nil)

		ctx := context.Background()
		result, err := service.GetByID(ctx, DefaultID)

		assert.NoError(t, err)
		assert.Equal(t, mockTodo, result)
	})

	t.Run("error when find by id", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On("FindById", mock.Anything, mock.AnythingOfType("string")).Return(nil, ErrDefault)

		ctx := context.Background()
		result, err := service.GetByID(ctx, DefaultID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestTodoCreate(t *testing.T) {
	t.Run("success when create", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On("Store", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(mockTodo, nil)

		ctx := context.Background()
		result, err := service.Create(ctx, &models.Todo{})

		assert.NoError(t, err)
		assert.Equal(t, mockTodo, result)
	})

	t.Run("error when create", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On("Store", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil, ErrDefault)

		ctx := context.Background()
		result, err := service.Create(ctx, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestTodoUpdate(t *testing.T) {
	t.Run("success when update", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On(
			"CountFindByID",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(10, nil)
		mockRepository.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(mockTodo, nil)

		ctx := context.Background()
		result, err := service.Update(ctx, DefaultID, &models.Todo{})

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("error when count find by id", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On(
			"CountFindByID",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(0, ErrDefault)
		mockRepository.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(nil, nil)

		ctx := context.Background()
		result, err := service.Update(ctx, DefaultID, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("error when update", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On(
			"CountFindByID",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(10, nil)
		mockRepository.On(
			"Update",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*models.Todo"),
		).Return(nil, ErrDefault)

		ctx := context.Background()
		result, err := service.Update(ctx, DefaultID, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestTodoDelete(t *testing.T) {
	t.Run("success when delete", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

		ctx := context.Background()
		err := service.Delete(ctx, DefaultID)

		assert.NoError(t, err)
	})

	t.Run("error when delete", func(t *testing.T) {
		mockRepository := new(mockRepositories.TodoRepository)
		service := services.NewTodoService(mockRepository)

		mockRepository.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(ErrDefault)

		ctx := context.Background()
		err := service.Delete(ctx, DefaultID)

		assert.Error(t, err)
	})
}
