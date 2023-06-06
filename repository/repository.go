package repository

import (
	"todoGin/model/entity"
)

type TodoRepository interface {
	GetAll() ([]entity.Todolist, error)
	GetByID(todoID int64) (*entity.Todolist, error)
	Create(title string) (*entity.Todolist, error)
	Update(todoID int64, updates map[string]interface{}) (*entity.Todolist, error)
	Delete(todoID int64) (int64, error)
	CreateUser(user *entity.User) error
	GetUserByUsername(username string) (*entity.User, error)
}
