package repository

import (
 "context"

 "github.com/taskflow/backend/internal/model"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
 // Create creates a new user in the database
 Create(ctx context.Context, user *model.User) error

 // GetByID retrieves a user by their ID
 GetByID(ctx context.Context, id string) (*model.User, error)

 // GetByEmail retrieves a user by their email address
 GetByEmail(ctx context.Context, email string) (*model.User, error)

 // Update updates an existing user
 Update(ctx context.Context, user *model.User) error

 // Delete deletes a user by their ID
 Delete(ctx context.Context, id string) error

 // List retrieves all users (for admin purposes)
 List(ctx context.Context, limit, offset int) ([]*model.User, error)
}
