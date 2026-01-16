package dto

import "time"

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
