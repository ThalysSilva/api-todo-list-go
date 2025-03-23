package services

import (
	"errors"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repositories"
)

type TodoService interface {
	CreateTodo(todo *models.Todo, userID int) error
	GetAllTodos() ([]models.Todo, error)
	UpdateTodo(todo *models.Todo, userID int) error
	DeleteTodo(id, userID int) error
}

type todoService struct {
	todoRepo repositories.TodoRepository
}

func NewTodoService(todoRepo repositories.TodoRepository) TodoService {
	return &todoService{todoRepo: todoRepo}
}

func (s *todoService) CreateTodo(todo *models.Todo, userID int) error {
	todo.UserID = userID
	return s.todoRepo.CreateTodo(todo)
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.todoRepo.GetAllTodos()
}

func (s *todoService) UpdateTodo(todo *models.Todo, userID int) error {
	existingTodo, err := s.todoRepo.GetTodoByID(todo.ID)
	if err != nil {
		return err
	}
	if existingTodo == nil {
		return errors.New("todo not found")
	}
	if existingTodo.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.todoRepo.UpdateTodo(todo)
}

func (s *todoService) DeleteTodo(id, userID int) error {
	todo, err := s.todoRepo.GetTodoByID(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return errors.New("todo not found")
	}
	if todo.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.todoRepo.DeleteTodo(id)
}