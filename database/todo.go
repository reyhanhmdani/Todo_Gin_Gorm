package database

import (
	"errors"
	"gorm.io/gorm"
	"todoGin/model/entity"
	"todoGin/repository"
)

// adaptop pattern
type TodoRepository struct {
	DB *gorm.DB
}

func NewTodoRepository(dbClient *gorm.DB) repository.TodoRepository {
	return &TodoRepository{
		DB: dbClient,
	}
}

func (t TodoRepository) GetAll() ([]entity.Todolist, error) {
	var todos []entity.Todolist

	result := t.DB.Find(&todos)
	return todos, result.Error
}

func (t TodoRepository) GetByID(todoID int64) (*entity.Todolist, error) {
	var todo entity.Todolist
	result := t.DB.Where("id = ?", todoID).First(&todo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &todo, result.Error
}

func (t TodoRepository) Create(title string) (*entity.Todolist, error) {
	todo := entity.Todolist{
		Title: title,
	}
	result := t.DB.Create(&todo)
	return &todo, result.Error
}

func (t TodoRepository) Update(todoID int64, updates map[string]interface{}) (*entity.Todolist, error) {
	var todo entity.Todolist
	result := t.DB.Model(&todo).Where("id = ?", todoID).Updates(updates)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &todo, result.Error
}

func (t TodoRepository) Delete(todoID int64) (int64, error) {
	todo := entity.Todolist{ID: todoID}
	result := t.DB.Delete(&todo)
	return result.RowsAffected, result.Error

}

func (t TodoRepository) CreateUser(user *entity.User) error {
	if err := t.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (t TodoRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := t.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil

}
