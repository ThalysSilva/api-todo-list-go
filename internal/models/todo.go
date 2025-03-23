package models

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
}