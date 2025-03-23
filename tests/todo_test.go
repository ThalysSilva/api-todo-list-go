package tests

import (
	"testing"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"
	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	todos map[int]*models.Todo
}

func (m *mockTodoRepository) CreateTodo(todo *models.Todo) error {
	todo.ID = len(m.todos) + 1
	m.todos[todo.ID] = todo
	return nil
}

func (m *mockTodoRepository) GetAllTodos() ([]models.Todo, error) {
	var todos []models.Todo
	for _, t := range m.todos {
		todos = append(todos, *t)
	}
	return todos, nil
}

func (m *mockTodoRepository) GetTodoByID(id int) (*models.Todo, error) {
	return m.todos[id], nil
}

func (m *mockTodoRepository) UpdateTodo(todo *models.Todo) error {
	m.todos[todo.ID] = todo
	return nil
}

func (m *mockTodoRepository) DeleteTodo(id int) error {
	delete(m.todos, id)
	return nil
}

func TestCreateTodo(t *testing.T) {
	repo := &mockTodoRepository{todos: make(map[int]*models.Todo)}
	service := services.NewTodoService(repo)
	todo := &models.Todo{Title: "Test", Description: "Test desc"}
	err := service.CreateTodo(todo, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, todo.UserID)
	assert.Equal(t, 1, todo.ID)
}