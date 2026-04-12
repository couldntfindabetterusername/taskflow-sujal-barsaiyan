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

type taskRepository struct {
 pool *pgxpool.Pool
}

// NewTaskRepository creates a new instance of TaskRepository
func NewTaskRepository(pool *pgxpool.Pool) repository.TaskRepository {
 return &taskRepository{
 pool: pool,
 }
}

// Create creates a new task in the database
func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
 // Generate UUID if not provided
 if task.ID == "" {
 task.ID = uuid.New().String()
 }

 query := `
 INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
 RETURNING created_at, updated_at
 `

 err := r.pool.QueryRow(ctx, query,
 task.ID,
 task.Title,
 task.Description,
 task.Status,
 task.Priority,
 task.ProjectID,
 task.AssigneeID,
 task.DueDate,
 ).Scan(&task.CreatedAt, &task.UpdatedAt)

 if err != nil {
 // Check for foreign key violations
 if strings.Contains(err.Error(), "foreign key") || strings.Contains(err.Error(), "violates foreign key constraint") {
 if strings.Contains(err.Error(), "project") {
 return fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 if strings.Contains(err.Error(), "assignee") || strings.Contains(err.Error(), "user") {
 return fmt.Errorf("%w: assignee not found", errors.ErrNotFound)
 }
 return fmt.Errorf("%w: referenced resource not found", errors.ErrNotFound)
 }
 return fmt.Errorf("failed to create task: %w", err)
 }

 return nil
}

// GetByID retrieves a task by its ID
func (r *taskRepository) GetByID(ctx context.Context, id string) (*model.Task, error) {
 query := `
 SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
 FROM tasks
 WHERE id = $1
 `

 var task model.Task
 err := r.pool.QueryRow(ctx, query, id).Scan(
 &task.ID,
 &task.Title,
 &task.Description,
 &task.Status,
 &task.Priority,
 &task.ProjectID,
 &task.AssigneeID,
 &task.DueDate,
 &task.CreatedAt,
 &task.UpdatedAt,
 )

 if err != nil {
 if err == pgx.ErrNoRows {
 return nil, fmt.Errorf("%w: task not found", errors.ErrNotFound)
 }
 return nil, fmt.Errorf("failed to get task by id: %w", err)
 }

 return &task, nil
}

// ListByProject retrieves all tasks for a specific project with optional filters
func (r *taskRepository) ListByProject(ctx context.Context, projectID string, filters *model.TaskFilters, limit, offset int) ([]*model.Task, error) {
 // Build query with dynamic filters
 query := `
 SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
 FROM tasks
 WHERE project_id = $1
 `
 args := []interface{}{projectID}
 argIndex := 2

 if filters != nil {
 if filters.Status != nil {
 query += fmt.Sprintf(" AND status = $%d", argIndex)
 args = append(args, *filters.Status)
 argIndex++
 }
 if filters.Priority != nil {
 query += fmt.Sprintf(" AND priority = $%d", argIndex)
 args = append(args, *filters.Priority)
 argIndex++
 }
 if filters.AssigneeID != nil {
 query += fmt.Sprintf(" AND assignee_id = $%d", argIndex)
 args = append(args, *filters.AssigneeID)
 argIndex++
 }
 }

 query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
 args = append(args, limit, offset)

 rows, err := r.pool.Query(ctx, query, args...)
 if err != nil {
 return nil, fmt.Errorf("failed to list tasks by project: %w", err)
 }
 defer rows.Close()

 var tasks []*model.Task
 for rows.Next() {
 var task model.Task
 err := rows.Scan(
 &task.ID,
 &task.Title,
 &task.Description,
 &task.Status,
 &task.Priority,
 &task.ProjectID,
 &task.AssigneeID,
 &task.DueDate,
 &task.CreatedAt,
 &task.UpdatedAt,
 )
 if err != nil {
 return nil, fmt.Errorf("failed to scan task: %w", err)
 }
 tasks = append(tasks, &task)
 }

 if err := rows.Err(); err != nil {
 return nil, fmt.Errorf("error iterating tasks: %w", err)
 }

 return tasks, nil
}

// ListByAssignee retrieves all tasks assigned to a specific user
func (r *taskRepository) ListByAssignee(ctx context.Context, assigneeID string, filters *model.TaskFilters, limit, offset int) ([]*model.Task, error) {
 query := `
 SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
 FROM tasks
 WHERE assignee_id = $1
 `
 args := []interface{}{assigneeID}
 argIndex := 2

 if filters != nil {
 if filters.Status != nil {
 query += fmt.Sprintf(" AND status = $%d", argIndex)
 args = append(args, *filters.Status)
 argIndex++
 }
 if filters.Priority != nil {
 query += fmt.Sprintf(" AND priority = $%d", argIndex)
 args = append(args, *filters.Priority)
 argIndex++
 }
 }

 query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
 args = append(args, limit, offset)

 rows, err := r.pool.Query(ctx, query, args...)
 if err != nil {
 return nil, fmt.Errorf("failed to list tasks by assignee: %w", err)
 }
 defer rows.Close()

 var tasks []*model.Task
 for rows.Next() {
 var task model.Task
 err := rows.Scan(
 &task.ID,
 &task.Title,
 &task.Description,
 &task.Status,
 &task.Priority,
 &task.ProjectID,
 &task.AssigneeID,
 &task.DueDate,
 &task.CreatedAt,
 &task.UpdatedAt,
 )
 if err != nil {
 return nil, fmt.Errorf("failed to scan task: %w", err)
 }
 tasks = append(tasks, &task)
 }

 if err := rows.Err(); err != nil {
 return nil, fmt.Errorf("error iterating tasks: %w", err)
 }

 return tasks, nil
}

// Update updates an existing task
func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
 query := `
 UPDATE tasks
 SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = NOW()
 WHERE id = $7
 RETURNING updated_at
 `

 err := r.pool.QueryRow(ctx, query,
 task.Title,
 task.Description,
 task.Status,
 task.Priority,
 task.AssigneeID,
 task.DueDate,
 task.ID,
 ).Scan(&task.UpdatedAt)

 if err != nil {
 if err == pgx.ErrNoRows {
 return fmt.Errorf("%w: task not found", errors.ErrNotFound)
 }
 // Check for foreign key violations
 if strings.Contains(err.Error(), "foreign key") || strings.Contains(err.Error(), "violates foreign key constraint") {
 return fmt.Errorf("%w: assignee not found", errors.ErrNotFound)
 }
 return fmt.Errorf("failed to update task: %w", err)
 }

 return nil
}

// Delete deletes a task by its ID
func (r *taskRepository) Delete(ctx context.Context, id string) error {
 query := `
 DELETE FROM tasks
 WHERE id = $1
 `

 result, err := r.pool.Exec(ctx, query, id)
 if err != nil {
 return fmt.Errorf("failed to delete task: %w", err)
 }

 if result.RowsAffected() == 0 {
 return fmt.Errorf("%w: task not found", errors.ErrNotFound)
 }

 return nil
}

// CountByProject returns the number of tasks in a project
func (r *taskRepository) CountByProject(ctx context.Context, projectID string) (int, error) {
 query := `
 SELECT COUNT(*) FROM tasks WHERE project_id = $1
 `

 var count int
 err := r.pool.QueryRow(ctx, query, projectID).Scan(&count)
 if err != nil {
 return 0, fmt.Errorf("failed to count tasks by project: %w", err)
 }

 return count, nil
}

// CountByStatus returns the number of tasks with a specific status in a project
func (r *taskRepository) CountByStatus(ctx context.Context, projectID string, status model.TaskStatus) (int, error) {
 query := `
 SELECT COUNT(*) FROM tasks WHERE project_id = $1 AND status = $2
 `

 var count int
 err := r.pool.QueryRow(ctx, query, projectID, status).Scan(&count)
 if err != nil {
 return 0, fmt.Errorf("failed to count tasks by status: %w", err)
 }

 return count, nil
}

// GetProjectStats returns task statistics for a project
func (r *taskRepository) GetProjectStats(ctx context.Context, projectID string) (*model.ProjectStats, error) {
 stats := &model.ProjectStats{
 ProjectID:  projectID,
 ByStatus:   make(map[string]int),
 ByAssignee: make(map[string]int),
 }

 // Get total count and completion percentage
 var totalCount, doneCount int
 countQuery := `
 SELECT
 COUNT(*) as total,
 COUNT(CASE WHEN status = 'done' THEN 1 END) as done
 FROM tasks
 WHERE project_id = $1
 `
 err := r.pool.QueryRow(ctx, countQuery, projectID).Scan(&totalCount, &doneCount)
 if err != nil {
 return nil, fmt.Errorf("failed to get task counts: %w", err)
 }

 stats.TotalTasks = totalCount
 if totalCount > 0 {
 stats.CompletionPercentage = float64(doneCount) / float64(totalCount) * 100
 }

 // Get counts by status
 statusQuery := `
 SELECT status, COUNT(*) as count
 FROM tasks
 WHERE project_id = $1
 GROUP BY status
 `
 rows, err := r.pool.Query(ctx, statusQuery, projectID)
 if err != nil {
 return nil, fmt.Errorf("failed to get status counts: %w", err)
 }
 defer rows.Close()

 for rows.Next() {
 var status string
 var count int
 if err := rows.Scan(&status, &count); err != nil {
 return nil, fmt.Errorf("failed to scan status count: %w", err)
 }
 stats.ByStatus[status] = count
 }

 // Get counts by assignee
 assigneeQuery := `
 SELECT
 COALESCE(assignee_id, 'unassigned') as assignee,
 COUNT(*) as count
 FROM tasks
 WHERE project_id = $1
 GROUP BY assignee_id
 `
 rows, err = r.pool.Query(ctx, assigneeQuery, projectID)
 if err != nil {
 return nil, fmt.Errorf("failed to get assignee counts: %w", err)
 }
 defer rows.Close()

 for rows.Next() {
 var assignee string
 var count int
 if err := rows.Scan(&assignee, &count); err != nil {
 return nil, fmt.Errorf("failed to scan assignee count: %w", err)
 }
 stats.ByAssignee[assignee] = count
 }

 return stats, nil
}
