package tests

import (
	"testing"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users map[string]*models.User
}

func (m *mockUserRepository) CreateUser(user *models.User) error {
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	return m.users[username], nil
}

func (m *mockUserRepository) UpdateRefreshToken(userID int, refreshToken string) error {
	for _, u := range m.users {
		if u.ID == userID {
			u.RefreshToken = refreshToken
			return nil
		}
	}
	return nil
}

func TestRegister(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]*models.User)}
	service := services.NewAuthService(repo)
	err := service.Register("testuser", "password123")
	assert.NoError(t, err)
	user, _ := repo.GetUserByUsername("testuser")
	assert.Equal(t, "testuser", user.Username)
}