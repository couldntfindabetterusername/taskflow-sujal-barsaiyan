package service

import (
 "context"
 "fmt"
 "log/slog"
 "time"

 "github.com/taskflow/backend/internal/errors"
 "github.com/taskflow/backend/internal/model"
 "github.com/taskflow/backend/internal/repository"
)

// TaskService handles task-related business logic
type TaskService struct {
 taskRepo    repository.TaskRepository
 projectRepo repository.ProjectRepository
 userRepo    repository.UserRepository
}

// NewTaskService creates a new TaskService instance
func NewTaskService(
 taskRepo repository.TaskRepository,
 projectRepo repository.ProjectRepository,
 userRepo repository.UserRepository,
) *TaskService {
 return &TaskService{
 taskRepo:    taskRepo,
 projectRepo: projectRepo,
 userRepo:    userRepo,
 }
}

// CreateTask creates a new task in a project
func (s *TaskService) CreateTask(ctx context.Context, projectID string, req model.TaskCreateRequest, creatorID string) (*model.TaskResponse, error) {
 // Validate input
 if err := req.Validate(); err != nil {
 return nil, errors.NewValidationError("task", err.Error())
 }

 // Validate due date if provided (must be in future)
 if req.DueDate != nil && req.DueDate.Before(time.Now()) {
 return nil, errors.NewValidationError("due_date", "due date must be in the future")
 }

 // Check if project exists
 project, err := s.projectRepo.GetByID(ctx, projectID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get project", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to get project: %w", err)
 }

 // Check if user has access to the project (owner or has tasks assigned)
 hasAccess, err := s.userHasProjectAccess(ctx, projectID, creatorID)
 if err != nil {
 return nil, err
 }
 if !hasAccess {
 return nil, fmt.Errorf("%w: you don't have access to this project", errors.ErrUnauthorized)
 }

 // Validate assignee if provided
 if req.AssigneeID != nil && *req.AssigneeID != "" {
 _, err := s.userRepo.GetByID(ctx, *req.AssigneeID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, errors.NewValidationError("assignee_id", "assignee not found")
 }
 slog.Error("Failed to validate assignee", "error", err, "assignee_id", *req.AssigneeID)
 return nil, fmt.Errorf("failed to validate assignee: %w", err)
 }
 }

 // Create task
 task := &model.Task{
 Title:       req.Title,
 Description: req.Description,
 Status:      req.Status,
 Priority:    req.Priority,
 ProjectID:   projectID,
 AssigneeID:  req.AssigneeID,
 DueDate:     req.DueDate,
 }

 if err := s.taskRepo.Create(ctx, task); err != nil {
 slog.Error("Failed to create task", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to create task: %w", err)
 }

 slog.Info("Task created successfully", "task_id", task.ID, "project_id", projectID, "creator_id", creatorID)

 // Build response with assignee details if available
 resp := task.ToResponse()
 if task.AssigneeID != nil {
 assignee, err := s.userRepo.GetByID(ctx, *task.AssigneeID)
 if err == nil {
 assigneeResp := assignee.ToResponse()
 resp.Assignee = &assigneeResp
 }
 }

 // Store creator ID for authorization purposes (not in DB, but tracked via project ownership)
 _ = project // Used for access check above

 return &resp, nil
}

// ListTasks retrieves tasks for a project with optional filters
func (s *TaskService) ListTasks(ctx context.Context, projectID string, userID string, filters *model.TaskFilters, limit, offset int) ([]model.TaskResponse, error) {
 // Check if project exists
 _, err := s.projectRepo.GetByID(ctx, projectID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get project", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to get project: %w", err)
 }

 // Check if user has access to the project
 hasAccess, err := s.userHasProjectAccess(ctx, projectID, userID)
 if err != nil {
 return nil, err
 }
 if !hasAccess {
 return nil, fmt.Errorf("%w: you don't have access to this project", errors.ErrUnauthorized)
 }

 if limit <= 0 {
 limit = 100 // Default limit
 }
 if limit > 1000 {
 limit = 1000 // Max limit
 }

 tasks, err := s.taskRepo.ListByProject(ctx, projectID, filters, limit, offset)
 if err != nil {
 slog.Error("Failed to list tasks", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to list tasks: %w", err)
 }

 responses := make([]model.TaskResponse, len(tasks))
 for i, task := range tasks {
 responses[i] = task.ToResponse()
 }

 return responses, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(ctx context.Context, taskID string, req model.TaskUpdateRequest, userID string) (*model.TaskResponse, error) {
 // Validate input
 if err := req.Validate(); err != nil {
 return nil, errors.NewValidationError("task", err.Error())
 }

 // Validate due date if provided (must be in future)
 if req.DueDate != nil && req.DueDate.Before(time.Now()) {
 return nil, errors.NewValidationError("due_date", "due date must be in the future")
 }

 // Get existing task
 task, err := s.taskRepo.GetByID(ctx, taskID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: task not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get task", "error", err, "task_id", taskID)
 return nil, fmt.Errorf("failed to get task: %w", err)
 }

 // Check if user has access to the project
 hasAccess, err := s.userHasProjectAccess(ctx, task.ProjectID, userID)
 if err != nil {
 return nil, err
 }
 if !hasAccess {
 return nil, fmt.Errorf("%w: you don't have access to this task", errors.ErrUnauthorized)
 }

 // Validate assignee if provided
 if req.AssigneeID != nil && *req.AssigneeID != "" {
 _, err := s.userRepo.GetByID(ctx, *req.AssigneeID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, errors.NewValidationError("assignee_id", "assignee not found")
 }
 slog.Error("Failed to validate assignee", "error", err, "assignee_id", *req.AssigneeID)
 return nil, fmt.Errorf("failed to validate assignee: %w", err)
 }
 }

 // Apply updates
 if req.Title != nil {
 if *req.Title == "" {
 return nil, errors.NewValidationError("title", "title cannot be empty")
 }
 task.Title = *req.Title
 }
 if req.Description != nil {
 task.Description = req.Description
 }
 if req.Status != nil {
 task.Status = *req.Status
 }
 if req.Priority != nil {
 task.Priority = *req.Priority
 }
 if req.AssigneeID != nil {
 // Allow setting assignee to null by passing empty string
 if *req.AssigneeID == "" {
 task.AssigneeID = nil
 } else {
 task.AssigneeID = req.AssigneeID
 }
 }
 if req.DueDate != nil {
 task.DueDate = req.DueDate
 }

 // Save updates
 if err := s.taskRepo.Update(ctx, task); err != nil {
 slog.Error("Failed to update task", "error", err, "task_id", taskID)
 return nil, fmt.Errorf("failed to update task: %w", err)
 }

 slog.Info("Task updated successfully", "task_id", taskID, "user_id", userID)

 // Build response with assignee details if available
 resp := task.ToResponse()
 if task.AssigneeID != nil {
 assignee, err := s.userRepo.GetByID(ctx, *task.AssigneeID)
 if err == nil {
 assigneeResp := assignee.ToResponse()
 resp.Assignee = &assigneeResp
 }
 }

 return &resp, nil
}

// DeleteTask deletes a task (project owner only)
func (s *TaskService) DeleteTask(ctx context.Context, taskID string, userID string) error {
 // Get existing task
 task, err := s.taskRepo.GetByID(ctx, taskID)
 if err != nil {
 if errors.IsNotFound(err) {
 return fmt.Errorf("%w: task not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get task", "error", err, "task_id", taskID)
 return fmt.Errorf("failed to get task: %w", err)
 }

 // Check if user is project owner (only project owner can delete tasks)
 isOwner, err := s.projectRepo.IsOwner(ctx, task.ProjectID, userID)
 if err != nil {
 slog.Error("Failed to check project ownership", "error", err, "project_id", task.ProjectID)
 return fmt.Errorf("failed to check authorization: %w", err)
 }

 if !isOwner {
 slog.Warn("Unauthorized task delete attempt", "task_id", taskID, "user_id", userID, "project_id", task.ProjectID)
 return fmt.Errorf("%w: only the project owner can delete tasks", errors.ErrUnauthorized)
 }

 // Delete task
 if err := s.taskRepo.Delete(ctx, taskID); err != nil {
 slog.Error("Failed to delete task", "error", err, "task_id", taskID)
 return fmt.Errorf("failed to delete task: %w", err)
 }

 slog.Info("Task deleted successfully", "task_id", taskID, "user_id", userID)
 return nil
}

// GetTask retrieves a single task by ID
func (s *TaskService) GetTask(ctx context.Context, taskID string, userID string) (*model.TaskResponse, error) {
 // Get task
 task, err := s.taskRepo.GetByID(ctx, taskID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: task not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get task", "error", err, "task_id", taskID)
 return nil, fmt.Errorf("failed to get task: %w", err)
 }

 // Check if user has access to the project
 hasAccess, err := s.userHasProjectAccess(ctx, task.ProjectID, userID)
 if err != nil {
 return nil, err
 }
 if !hasAccess {
 return nil, fmt.Errorf("%w: you don't have access to this task", errors.ErrUnauthorized)
 }

 // Build response with assignee details if available
 resp := task.ToResponse()
 if task.AssigneeID != nil {
 assignee, err := s.userRepo.GetByID(ctx, *task.AssigneeID)
 if err == nil {
 assigneeResp := assignee.ToResponse()
 resp.Assignee = &assigneeResp
 }
 }

 return &resp, nil
}

// userHasProjectAccess checks if a user has access to a project (owner or has assigned tasks)
func (s *TaskService) userHasProjectAccess(ctx context.Context, projectID, userID string) (bool, error) {
 // Check if user is owner
 isOwner, err := s.projectRepo.IsOwner(ctx, projectID, userID)
 if err != nil {
 return false, fmt.Errorf("failed to check project ownership: %w", err)
 }
 if isOwner {
 return true, nil
 }

 // Check if user has any tasks assigned in this project
 tasks, err := s.taskRepo.ListByProject(ctx, projectID, &model.TaskFilters{AssigneeID: &userID}, 1, 0)
 if err != nil {
 return false, fmt.Errorf("failed to check task assignments: %w", err)
 }

 return len(tasks) > 0, nil
}
