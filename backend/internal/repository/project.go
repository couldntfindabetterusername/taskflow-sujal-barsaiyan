package repository

import (
 "context"

 "github.com/taskflow/backend/internal/model"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
 // Create creates a new project in the database
 Create(ctx context.Context, project *model.Project) error

 // GetByID retrieves a project by its ID
 GetByID(ctx context.Context, id string) (*model.Project, error)

 // ListByUser retrieves all projects owned by a specific user
 ListByUser(ctx context.Context, userID string, limit, offset int) ([]*model.Project, error)

 // Update updates an existing project
 Update(ctx context.Context, project *model.Project) error

 // Delete deletes a project by its ID
 Delete(ctx context.Context, id string) error

 // List retrieves all projects (for admin purposes)
 List(ctx context.Context, limit, offset int) ([]*model.Project, error)

 // IsOwner checks if a user is the owner of a project
 IsOwner(ctx context.Context, projectID, userID string) (bool, error)
}
