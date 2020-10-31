package repository

import (
	"errors"
	"os"

	"../models"

	"github.com/juju/mgosession"
	"gopkg.in/mgo.v2/bson"
)

// TodoRepository represent the todo repository contract
type TodoRepository interface {
	FindAll(keyword string, limit int, offset int) ([]*models.Todo, error)
	CountFindAll(keyword string) (int, error)
	FindById(id string) (*models.Todo, error)
	CountFindByID(id string) (int, error)
	Store(value interface{}) (*models.Todo, error)
	Update(id string, value interface{}) error
	Delete(id string) error
}

type mongoTodoRepository struct {
	pool *mgosession.Pool
}

// NewMongoTodoRepository will create an object that represent the TodoRepository interface
func NewMongoTodoRepository(Pool *mgosession.Pool) TodoRepository {
	return &mongoTodoRepository{Pool}
}

// FindAll - find all todo
func (m *mongoTodoRepository) FindAll(keyword string, limit int, offset int) ([]*models.Todo, error) {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")
	results := []*models.Todo{}

	err := c.Find(bson.M{"title": bson.M{"$regex": keyword, "$options": "i"}}).Sort("-createdAt").
		Limit(limit).Skip(offset).All(&results)
	if err != nil {
		return results, err
	}

	return results, nil
}

// CountFindAll - count find all todo
func (m *mongoTodoRepository) CountFindAll(keyword string) (int, error) {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")

	total, err := c.Find(bson.M{"title": bson.M{"$regex": keyword, "$options": "i"}}).Count()
	if err != nil {
		return total, err
	}

	return total, nil
}

// FindById - find todo by id
func (m *mongoTodoRepository) FindById(id string) (*models.Todo, error) {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")
	result := &models.Todo{}

	// Validate string if hex
	if !bson.IsObjectIdHex(id) {
		return result, errors.New("not found")
	}

	// Find data
	err := c.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// CountFindByID - find count todo by id
func (m *mongoTodoRepository) CountFindByID(id string) (int, error) {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")

	// Validate string if hex
	if !bson.IsObjectIdHex(id) {
		return 0, errors.New("not found")
	}

	// Find data
	total, err := c.FindId(bson.ObjectIdHex(id)).Count()
	if err != nil {
		return total, err
	}

	return total, nil
}

// Store - store todo
func (m *mongoTodoRepository) Store(value interface{}) (*models.Todo, error) {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")
	result := &models.Todo{}

	// Insert data
	info, err := c.UpsertId(bson.NewObjectId(), value)
	if err != nil {
		return result, err
	}

	// Find data
	err = c.FindId(info.UpsertedId).One(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Update - update todo
func (m *mongoTodoRepository) Update(id string, value interface{}) error {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")

	// Validate string if hex
	if !bson.IsObjectIdHex(id) {
		return errors.New("not valid")
	}

	// Find data
	err := c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, value)
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete todo
func (m *mongoTodoRepository) Delete(id string) error {
	session := m.pool.Session(nil)
	c := session.DB(os.Getenv("DB_NAME")).C("todo")

	// Validate string if hex
	if !bson.IsObjectIdHex(id) {
		return errors.New("not valid")
	}

	// Find data
	err := c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		return err
	}

	return nil
}
