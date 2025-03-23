package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username" validate:"required"`
	PasswordHash string `json:"-"`
	RefreshToken string `json:"-"`
}