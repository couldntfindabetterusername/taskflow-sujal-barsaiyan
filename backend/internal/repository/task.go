package repository

import (
 "context"

 "github.com/taskflow/backend/internal/model"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
 // Create creates a new task in the database
 Create(ctx context.Context, task *model.Task) error

 // GetByID retrieves a task by its ID
 GetByID(ctx context.Context, id string) (*model.Task, error)

 // ListByProject retrieves all tasks for a specific project with optional filters
 ListByProject(ctx context.Context, projectID string, filters *model.TaskFilters, limit, offset int) ([]*model.Task, error)

 // ListByAssignee retrieves all tasks assigned to a specific user
 ListByAssignee(ctx context.Context, assigneeID string, filters *model.TaskFilters, limit, offset int) ([]*model.Task, error)

 // Update updates an existing task
 Update(ctx context.Context, task *model.Task) error

 // Delete deletes a task by its ID
 Delete(ctx context.Context, id string) error

 // CountByProject returns the number of tasks in a project
 CountByProject(ctx context.Context, projectID string) (int, error)

 // CountByStatus returns the number of tasks with a specific status in a project
 CountByStatus(ctx context.Context, projectID string, status model.TaskStatus) (int, error)
}
