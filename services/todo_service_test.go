package services_test

import (
	"errors"
	"testing"

	mockRepositories "go-clean-architecture/mocks/repository"
	"go-clean-architecture/models"
	services "go-clean-architecture/services"
	utils "go-clean-architecture/utils"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

// TestGetAllSuccess - testing GetAll
func TestGetAllSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)
	mockListTodo := make([]*models.Todo, 0)
	mockListTodo = append(mockListTodo, &mockTodo)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(mockListTodo, nil)
	mockRepository.On("CountFindAll", mock.Anything).Return(10, nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	_, _, err = todoService.GetAll("keyword", 10, 0)
	assert.NoError(t, err)
}

// TestGetAllDatabaseError - testing GetAll
func TestGetAllFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)
	mockListTodo := make([]*models.Todo, 0)
	mockListTodo = append(mockListTodo, &mockTodo)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(mockListTodo, nil)
	mockRepository.On("CountFindAll", mock.Anything).Return(10, nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	_, _, err = todoService.GetAll("keyword", 10, 0)
	assert.NoError(t, err)

	mockRepository = new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(mockListTodo, errors.New(""))
	mockRepository.On("CountFindAll", mock.Anything).Return(10, nil)

	// Mock service in handler
	todoService = services.NewTodoService(mockRepository)

	_, _, err = todoService.GetAll("keyword", 10, 0)
	assert.NoError(t, err)

	mockRepository = new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(mockListTodo, nil)
	mockRepository.On("CountFindAll", mock.Anything).Return(10, errors.New(""))

	// Mock service in handler
	todoService = services.NewTodoService(mockRepository)

	_, _, err = todoService.GetAll("keyword", 10, 0)
	assert.NoError(t, err)
}

// TestGetByIDSuccess - testing GetByID
func TestGetByIDSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("FindById", mock.AnythingOfType("string")).Return(&mockTodo, nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	_, err = todoService.GetByID("5f9e9dc1c5e4cd34e6ac0bb7")
	assert.NoError(t, err)
}

// TestGetByIDFailedDatabaseError - testing GetByID
func TestGetByIDFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("FindById", mock.AnythingOfType("string")).Return(&mockTodo, errors.New(""))

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	todoService.GetByID("5f9e9dc1c5e4cd34e6ac0bb7")
}

// TestCreateSuccess - testing Create
func TestCreateSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("Store", mock.AnythingOfType("primitive.M")).Return(&mockTodo, nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	timeNow := utils.GetTimeNow()

	_, err = todoService.Create(bson.M{
		"title":       "lorem ipsum",
		"description": "lorem ipsum",
		"createdAt":   timeNow,
		"updatedAt":   timeNow,
		"deletedAt":   timeNow,
	})
	assert.NoError(t, err)
}

// TestCreateFailedDatabaseError - testing Create
func TestCreateFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("Store", mock.AnythingOfType("primitive.M")).Return(&mockTodo, errors.New(""))

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	timeNow := utils.GetTimeNow()

	todoService.Create(bson.M{
		"title":       "lorem ipsum",
		"description": "lorem ipsum",
		"createdAt":   timeNow,
		"updatedAt":   timeNow,
		"deletedAt":   timeNow,
	})
}

// TestUpdateSuccess - testing Update
func TestUpdateSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("CountFindByID", mock.AnythingOfType("string")).Return(1, nil)
	mockRepository.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("primitive.D")).Return(&mockTodo, nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	timeNow := utils.GetTimeNow()

	todoService.Update("5f9e9dc1c5e4cd34e6ac0bb7", bson.D{
		{Key: "title", Value: "lorem ipsum"},
		{Key: "description", Value: "lorem ipsum"},
		{Key: "updatedAt", Value: timeNow},
	})
}

// TestUpdateFailedDatabaseError - testing Update
func TestUpdateFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("CountFindByID", mock.AnythingOfType("string")).Return(1, errors.New(""))
	mockRepository.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("primitive.D")).Return(&mockTodo, nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	timeNow := utils.GetTimeNow()

	todoService.Update("5f9e9dc1c5e4cd34e6ac0bb7", bson.D{
		{Key: "title", Value: "lorem ipsum"},
		{Key: "description", Value: "lorem ipsum"},
		{Key: "updatedAt", Value: timeNow},
	})

	mockRepository = new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("CountFindByID", mock.AnythingOfType("string")).Return(1, nil)
	mockRepository.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("primitive.D")).Return(&mockTodo, errors.New(""))

	// Mock service in handler
	todoService = services.NewTodoService(mockRepository)

	timeNow = utils.GetTimeNow()

	todoService.Update("5f9e9dc1c5e4cd34e6ac0bb7", bson.D{
		{Key: "title", Value: "lorem ipsum"},
		{Key: "description", Value: "lorem ipsum"},
		{Key: "updatedAt", Value: timeNow},
	})
}

// TestDeleteSuccess - testing Delete
func TestDeleteSuccess(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("Delete", mock.AnythingOfType("string")).Return(nil)

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	todoService.Delete("5f9e9dc1c5e4cd34e6ac0bb7")
}

// TestDeleteFailedDatabaseError - testing Delete
func TestDeleteFailedDatabaseError(t *testing.T) {
	var mockTodo models.Todo
	err := faker.FakeData(&mockTodo)
	assert.NoError(t, err)

	mockRepository := new(mockRepositories.TodoRepository)

	// Mock GetAll service
	mockRepository.On("Delete", mock.AnythingOfType("string")).Return(errors.New(""))

	// Mock service in handler
	todoService := services.NewTodoService(mockRepository)

	todoService.Delete("5f9e9dc1c5e4cd34e6ac0bb7")
}
