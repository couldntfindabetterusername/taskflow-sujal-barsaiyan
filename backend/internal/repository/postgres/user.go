package postgres

import (
 "context"
 "fmt"
 "strings"

 "github.com/google/uuid"
 "github.com/jackc/pgx/v5"
 "github.com/jackc/pgx/v5/pgxpool"
 "github.com/taskflow/backend/internal/errors"
 "github.com/taskflow/backend/internal/model"
 "github.com/taskflow/backend/internal/repository"
)

type userRepository struct {
 pool *pgxpool.Pool
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
 return &userRepository{
 pool: pool,
 }
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
 // Generate UUID if not provided
 if user.ID == "" {
 user.ID = uuid.New().String()
 }

 query := `
 INSERT INTO users (id, name, email, password_hash, created_at)
 VALUES ($1, $2, $3, $4, NOW())
 RETURNING created_at
 `

 err := r.pool.QueryRow(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash).Scan(&user.CreatedAt)
 if err != nil {
 // Check for duplicate email error
 if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
 return fmt.Errorf("%w: email already exists", errors.ErrDuplicate)
 }
 return fmt.Errorf("failed to create user: %w", err)
 }

 return nil
}

// GetByID retrieves a user by their ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
 query := `
 SELECT id, name, email, password_hash, created_at
 FROM users
 WHERE id = $1
 `

 var user model.User
 err := r.pool.QueryRow(ctx, query, id).Scan(
 &user.ID,
 &user.Name,
 &user.Email,
 &user.PasswordHash,
 &user.CreatedAt,
 )

 if err != nil {
 if err == pgx.ErrNoRows {
 return nil, fmt.Errorf("%w: user not found", errors.ErrNotFound)
 }
 return nil, fmt.Errorf("failed to get user by id: %w", err)
 }

 return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
 query := `
 SELECT id, name, email, password_hash, created_at
 FROM users
 WHERE email = $1
 `

 var user model.User
 err := r.pool.QueryRow(ctx, query, email).Scan(
 &user.ID,
 &user.Name,
 &user.Email,
 &user.PasswordHash,
 &user.CreatedAt,
 )

 if err != nil {
 if err == pgx.ErrNoRows {
 return nil, fmt.Errorf("%w: user not found", errors.ErrNotFound)
 }
 return nil, fmt.Errorf("failed to get user by email: %w", err)
 }

 return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
 query := `
 UPDATE users
 SET name = $1, email = $2
 WHERE id = $3
 `

 result, err := r.pool.Exec(ctx, query, user.Name, user.Email, user.ID)
 if err != nil {
 // Check for duplicate email error
 if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
 return fmt.Errorf("%w: email already exists", errors.ErrDuplicate)
 }
 return fmt.Errorf("failed to update user: %w", err)
 }

 if result.RowsAffected() == 0 {
 return fmt.Errorf("%w: user not found", errors.ErrNotFound)
 }

 return nil
}

// Delete deletes a user by their ID
func (r *userRepository) Delete(ctx context.Context, id string) error {
 query := `
 DELETE FROM users
 WHERE id = $1
 `

 result, err := r.pool.Exec(ctx, query, id)
 if err != nil {
 return fmt.Errorf("failed to delete user: %w", err)
 }

 if result.RowsAffected() == 0 {
 return fmt.Errorf("%w: user not found", errors.ErrNotFound)
 }

 return nil
}

// List retrieves all users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
 query := `
 SELECT id, name, email, password_hash, created_at
 FROM users
 ORDER BY created_at DESC
 LIMIT $1 OFFSET $2
 `

 rows, err := r.pool.Query(ctx, query, limit, offset)
 if err != nil {
 return nil, fmt.Errorf("failed to list users: %w", err)
 }
 defer rows.Close()

 var users []*model.User
 for rows.Next() {
 var user model.User
 err := rows.Scan(
 &user.ID,
 &user.Name,
 &user.Email,
 &user.PasswordHash,
 &user.CreatedAt,
 )
 if err != nil {
 return nil, fmt.Errorf("failed to scan user: %w", err)
 }
 users = append(users, &user)
 }

 if err := rows.Err(); err != nil {
 return nil, fmt.Errorf("error iterating users: %w", err)
 }

 return users, nil
}
