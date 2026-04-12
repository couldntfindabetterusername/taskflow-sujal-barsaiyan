package model

import "time"

// User represents a user in the system
type User struct {
 ID           string    `json:"id"`
 Name         string    `json:"name"`
 Email        string    `json:"email"`
 PasswordHash string    `json:"-"` // Never expose password hash in JSON
 CreatedAt    time.Time `json:"created_at"`
}

// UserCreateRequest represents the request body for creating a new user
type UserCreateRequest struct {
 Name     string `json:"name"`
 Email    string `json:"email"`
 Password string `json:"password"`
}

// UserUpdateRequest represents the request body for updating a user
type UserUpdateRequest struct {
 Name  *string `json:"name,omitempty"`
 Email *string `json:"email,omitempty"`
}

// UserResponse represents a user in API responses (without sensitive data)
type UserResponse struct {
 ID        string    `json:"id"`
 Name      string    `json:"name"`
 Email     string    `json:"email"`
 CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
 return UserResponse{
 ID:        u.ID,
 Name:      u.Name,
 Email:     u.Email,
 CreatedAt: u.CreatedAt,
 }
}
