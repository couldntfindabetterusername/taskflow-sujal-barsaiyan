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

type projectRepository struct {
 pool *pgxpool.Pool
}

// NewProjectRepository creates a new instance of ProjectRepository
func NewProjectRepository(pool *pgxpool.Pool) repository.ProjectRepository {
 return &projectRepository{
 pool: pool,
 }
}

// Create creates a new project in the database
func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
 // Generate UUID if not provided
 if project.ID == "" {
 project.ID = uuid.New().String()
 }

 query := `
 INSERT INTO projects (id, name, description, owner_id, created_at)
 VALUES ($1, $2, $3, $4, NOW())
 RETURNING created_at
 `

 err := r.pool.QueryRow(ctx, query, project.ID, project.Name, project.Description, project.OwnerID).Scan(&project.CreatedAt)
 if err != nil {
 // Check for foreign key violation (invalid owner_id)
 if strings.Contains(err.Error(), "foreign key") || strings.Contains(err.Error(), "violates foreign key constraint") {
 return fmt.Errorf("%w: owner not found", errors.ErrNotFound)
 }
 return fmt.Errorf("failed to create project: %w", err)
 }

 return nil
}

// GetByID retrieves a project by its ID
func (r *projectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
 query := `
 SELECT id, name, description, owner_id, created_at
 FROM projects
 WHERE id = $1
 `

 var project model.Project
 err := r.pool.QueryRow(ctx, query, id).Scan(
 &project.ID,
 &project.Name,
 &project.Description,
 &project.OwnerID,
 &project.CreatedAt,
 )

 if err != nil {
 if err == pgx.ErrNoRows {
 return nil, fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 return nil, fmt.Errorf("failed to get project by id: %w", err)
 }

 return &project, nil
}

// ListByUser retrieves all projects owned by a specific user or where user has tasks
func (r *projectRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*model.Project, error) {
 // Get projects owned by user OR projects where user has assigned tasks
 query := `
 SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at
 FROM projects p
 LEFT JOIN tasks t ON t.project_id = p.id AND t.assignee_id = $1
 WHERE p.owner_id = $1 OR t.assignee_id = $1
 ORDER BY p.created_at DESC
 LIMIT $2 OFFSET $3
 `

 rows, err := r.pool.Query(ctx, query, userID, limit, offset)
 if err != nil {
 return nil, fmt.Errorf("failed to list projects by user: %w", err)
 }
 defer rows.Close()

 var projects []*model.Project
 for rows.Next() {
 var project model.Project
 err := rows.Scan(
 &project.ID,
 &project.Name,
 &project.Description,
 &project.OwnerID,
 &project.CreatedAt,
 )
 if err != nil {
 return nil, fmt.Errorf("failed to scan project: %w", err)
 }
 projects = append(projects, &project)
 }

 if err := rows.Err(); err != nil {
 return nil, fmt.Errorf("error iterating projects: %w", err)
 }

 return projects, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
 query := `
 UPDATE projects
 SET name = $1, description = $2
 WHERE id = $3
 `

 result, err := r.pool.Exec(ctx, query, project.Name, project.Description, project.ID)
 if err != nil {
 return fmt.Errorf("failed to update project: %w", err)
 }

 if result.RowsAffected() == 0 {
 return fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }

 return nil
}

// Delete deletes a project by its ID (tasks are cascade deleted via foreign key)
func (r *projectRepository) Delete(ctx context.Context, id string) error {
 query := `
 DELETE FROM projects
 WHERE id = $1
 `

 result, err := r.pool.Exec(ctx, query, id)
 if err != nil {
 return fmt.Errorf("failed to delete project: %w", err)
 }

 if result.RowsAffected() == 0 {
 return fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }

 return nil
}

// List retrieves all projects (for admin purposes)
func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]*model.Project, error) {
 query := `
 SELECT id, name, description, owner_id, created_at
 FROM projects
 ORDER BY created_at DESC
 LIMIT $1 OFFSET $2
 `

 rows, err := r.pool.Query(ctx, query, limit, offset)
 if err != nil {
 return nil, fmt.Errorf("failed to list projects: %w", err)
 }
 defer rows.Close()

 var projects []*model.Project
 for rows.Next() {
 var project model.Project
 err := rows.Scan(
 &project.ID,
 &project.Name,
 &project.Description,
 &project.OwnerID,
 &project.CreatedAt,
 )
 if err != nil {
 return nil, fmt.Errorf("failed to scan project: %w", err)
 }
 projects = append(projects, &project)
 }

 if err := rows.Err(); err != nil {
 return nil, fmt.Errorf("error iterating projects: %w", err)
 }

 return projects, nil
}

// IsOwner checks if a user is the owner of a project
func (r *projectRepository) IsOwner(ctx context.Context, projectID, userID string) (bool, error) {
 query := `
 SELECT EXISTS(
 SELECT 1 FROM projects
 WHERE id = $1 AND owner_id = $2
 )
 `

 var exists bool
 err := r.pool.QueryRow(ctx, query, projectID, userID).Scan(&exists)
 if err != nil {
 return false, fmt.Errorf("failed to check project ownership: %w", err)
 }

 return exists, nil
}
