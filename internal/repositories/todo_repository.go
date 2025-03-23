package repositories

import (
	"database/sql"
	"todo-list-api/internal/models"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
	GetAllTodos() ([]models.Todo, error)
	GetTodoByID(id int) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id int) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) CreateTodo(todo *models.Todo) error {
	query := "INSERT INTO todos (title, description, user_id) VALUES ($1, $2, $3) RETURNING id"
	return r.db.QueryRow(query, todo.Title, todo.Description, todo.UserID).Scan(&todo.ID)
}

func (r *todoRepository) GetAllTodos() ([]models.Todo, error) {
	query := `
		SELECT t.id, t.title, t.description, t.user_id, u.username
		FROM todos t
		JOIN users u ON t.user_id = u.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.UserID, &todo.Username); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *todoRepository) GetTodoByID(id int) (*models.Todo, error) {
	todo := &models.Todo{}
	query := "SELECT id, title, description, user_id FROM todos WHERE id = $1"
	err := r.db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.UserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepository) UpdateTodo(todo *models.Todo) error {
	query := "UPDATE todos SET title = $1, description = $2 WHERE id = $3"
	_, err := r.db.Exec(query, todo.Title, todo.Description, todo.ID)
	return err
}

func (r *todoRepository) DeleteTodo(id int) error {
	query := "DELETE FROM todos WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}