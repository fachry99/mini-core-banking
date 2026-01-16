package dto

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
