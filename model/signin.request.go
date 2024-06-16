package model

import "github.com/google/uuid"

type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type SigninOK struct {
	Username string `json:"username" binding:"required,email"`
	Jwt      string `json:"jwt"`
	Id       string `json:"id"`
}

type WhoamiResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
	Token    string    `json:"token"`
}

type EditUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type IsUsernameExistRequest struct {
	Username *string `json:"username,omitempty"`
}
