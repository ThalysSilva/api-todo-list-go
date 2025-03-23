package repositories

import (
	"database/sql"
	"todo-list-api/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	UpdateRefreshToken(userID int, refreshToken string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	return r.db.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID)
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, password_hash, refresh_token FROM users WHERE username = $1"
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.RefreshToken)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateRefreshToken(userID int, refreshToken string) error {
	query := "UPDATE users SET refresh_token = $1 WHERE id = $2"
	_, err := r.db.Exec(query, refreshToken, userID)
	return err
}