package services

import (
	"../models"
	"../repository"
)

// TodoService represent the instance service
type TodoService interface {
	GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error)
	GetByID(id string) (*models.Todo, error)
	CountGetByID(id string) (int, error)
	Create(value interface{}) (*models.Todo, error)
	Update(id string, value interface{}) error
	Delete(id string) error
}

type instanceService struct {
	instanceRepo repository.TodoRepository
}

// NewTodoService will create new an TodoService object representation of TodoService interface
func NewTodoService(a repository.TodoRepository) TodoService {
	return &instanceService{
		instanceRepo: a,
	}
}

// GetAll - get all instance service
func (a *instanceService) GetAll(keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	res, err := a.instanceRepo.FindAll(keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := a.instanceRepo.CountFindAll(keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// CountGetByID - count get instance by id service
func (a *instanceService) CountGetByID(id string) (int, error) {
	res, err := a.instanceRepo.CountFindByID(id)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetByID - get instance by id service
func (a *instanceService) GetByID(id string) (*models.Todo, error) {
	res, err := a.instanceRepo.FindById(id)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Create - creating instance service
func (a *instanceService) Create(value interface{}) (*models.Todo, error) {
	res, err := a.instanceRepo.Store(value)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Update - update instance service
func (a *instanceService) Update(id string, value interface{}) error {
	err := a.instanceRepo.Update(id, value)
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete instance service
func (a *instanceService) Delete(id string) error {
	err := a.instanceRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
