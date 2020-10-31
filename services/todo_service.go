package services

import (
	"../models"
	"../repository"
)

// TodoService represent the todo service
type TodoService interface {
	GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error)
	GetByID(id string) (*models.Todo, error)
	CountGetByID(id string) (int, error)
	Create(value interface{}) (*models.Todo, error)
	Update(id string, value interface{}) error
	Delete(id string) error
}

type todoService struct {
	todoRepo repository.TodoRepository
}

// NewTodoService will create new an TodoService object representation of TodoService interface
func NewTodoService(a repository.TodoRepository) TodoService {
	return &todoService{
		todoRepo: a,
	}
}

// GetAll - get all todo service
func (a *todoService) GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	res, err := a.todoRepo.FindAll(keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := a.todoRepo.CountFindAll(keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// CountGetByID - count get todo by id service
func (a *todoService) CountGetByID(id string) (int, error) {
	res, err := a.todoRepo.CountFindByID(id)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetByID - get todo by id service
func (a *todoService) GetByID(id string) (*models.Todo, error) {
	res, err := a.todoRepo.FindById(id)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Create - creating todo service
func (a *todoService) Create(value interface{}) (*models.Todo, error) {
	res, err := a.todoRepo.Store(value)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Update - update todo service
func (a *todoService) Update(id string, value interface{}) error {
	err := a.todoRepo.Update(id, value)
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete todo service
func (a *todoService) Delete(id string) error {
	err := a.todoRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
