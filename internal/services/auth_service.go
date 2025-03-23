package services

import (
	"errors"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repositories"
	"todo-list-api/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
	}
	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(username, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid password")
	}
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}
	if err := s.userRepo.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	_, err := auth.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}
	user, err := s.userRepo.GetUserByUsername("")
	if err != nil {
		return "", err
	}
	if user == nil || user.RefreshToken != refreshToken {
		return "", errors.New("invalid refresh token")
	}
	newAccessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}
	newRefreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", err
	}
	if err := s.userRepo.UpdateRefreshToken(user.ID, newRefreshToken); err != nil {
		return "", err
	}
	return newAccessToken, nil
}