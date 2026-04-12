package service

import (
 "context"
 "fmt"
 "log/slog"

 "github.com/taskflow/backend/internal/errors"
 "github.com/taskflow/backend/internal/model"
 "github.com/taskflow/backend/internal/repository"
)

// ProjectService handles project-related business logic
type ProjectService struct {
 projectRepo repository.ProjectRepository
 taskRepo    repository.TaskRepository
 userRepo    repository.UserRepository
}

// NewProjectService creates a new ProjectService instance
func NewProjectService(
 projectRepo repository.ProjectRepository,
 taskRepo repository.TaskRepository,
 userRepo repository.UserRepository,
) *ProjectService {
 return &ProjectService{
 projectRepo: projectRepo,
 taskRepo:    taskRepo,
 userRepo:    userRepo,
 }
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, req model.ProjectCreateRequest, ownerID string) (*model.ProjectResponse, error) {
 // Validate input
 if err := s.validateCreateRequest(req); err != nil {
 return nil, err
 }

 // Create project
 project := &model.Project{
 Name:        req.Name,
 Description: req.Description,
 OwnerID:     ownerID,
 }

 if err := s.projectRepo.Create(ctx, project); err != nil {
 slog.Error("Failed to create project", "error", err, "owner_id", ownerID)
 return nil, fmt.Errorf("failed to create project: %w", err)
 }

 slog.Info("Project created successfully", "project_id", project.ID, "owner_id", ownerID)

 // Get owner details for response
 owner, err := s.userRepo.GetByID(ctx, ownerID)
 if err != nil {
 slog.Warn("Failed to get owner details", "error", err, "owner_id", ownerID)
 // Return without owner details
 resp := project.ToResponse()
 return &resp, nil
 }

 resp := project.ToResponse()
 ownerResp := owner.ToResponse()
 resp.Owner = &ownerResp
 return &resp, nil
}

// GetProject retrieves a project by ID with tasks
func (s *ProjectService) GetProject(ctx context.Context, projectID string, userID string) (*ProjectDetailResponse, error) {
 // Get project
 project, err := s.projectRepo.GetByID(ctx, projectID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get project", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to get project: %w", err)
 }

 // Check if user has access (owner or has tasks assigned)
 hasAccess, err := s.userHasAccess(ctx, projectID, userID)
 if err != nil {
 return nil, err
 }
 if !hasAccess {
 return nil, fmt.Errorf("%w: you don't have access to this project", errors.ErrUnauthorized)
 }

 // Get tasks for this project
 tasks, err := s.taskRepo.ListByProject(ctx, projectID, nil, 1000, 0)
 if err != nil {
 slog.Warn("Failed to get project tasks", "error", err, "project_id", projectID)
 tasks = []*model.Task{} // Return empty slice if tasks fail
 }

 // Get owner details
 owner, err := s.userRepo.GetByID(ctx, project.OwnerID)
 if err != nil {
 slog.Warn("Failed to get owner details", "error", err, "owner_id", project.OwnerID)
 }

 resp := project.ToResponse()
 if owner != nil {
 ownerResp := owner.ToResponse()
 resp.Owner = &ownerResp
 }

 // Convert tasks to response format
 taskResponses := make([]model.TaskResponse, len(tasks))
 for i, task := range tasks {
 taskResponses[i] = task.ToResponse()
 }

 return &ProjectDetailResponse{
 Project: resp,
 Tasks:   taskResponses,
 }, nil
}

// ListProjects retrieves all projects accessible by a user
func (s *ProjectService) ListProjects(ctx context.Context, userID string, limit, offset int) ([]model.ProjectResponse, error) {
 if limit <= 0 {
 limit = 50 // Default limit
 }
 if limit > 100 {
 limit = 100 // Max limit
 }

 projects, err := s.projectRepo.ListByUser(ctx, userID, limit, offset)
 if err != nil {
 slog.Error("Failed to list projects", "error", err, "user_id", userID)
 return nil, fmt.Errorf("failed to list projects: %w", err)
 }

 responses := make([]model.ProjectResponse, len(projects))
 for i, project := range projects {
 responses[i] = project.ToResponse()
 }

 return responses, nil
}

// UpdateProject updates a project (owner only)
func (s *ProjectService) UpdateProject(ctx context.Context, projectID string, req model.ProjectUpdateRequest, userID string) (*model.ProjectResponse, error) {
 // Get existing project
 project, err := s.projectRepo.GetByID(ctx, projectID)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get project for update", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to get project: %w", err)
 }

 // Check ownership
 if project.OwnerID != userID {
 slog.Warn("Unauthorized project update attempt", "project_id", projectID, "user_id", userID, "owner_id", project.OwnerID)
 return nil, fmt.Errorf("%w: only the project owner can update this project", errors.ErrUnauthorized)
 }

 // Apply updates
 if req.Name != nil {
 if *req.Name == "" {
 return nil, errors.NewValidationError("name", "name cannot be empty")
 }
 project.Name = *req.Name
 }
 if req.Description != nil {
 project.Description = req.Description
 }

 // Save updates
 if err := s.projectRepo.Update(ctx, project); err != nil {
 slog.Error("Failed to update project", "error", err, "project_id", projectID)
 return nil, fmt.Errorf("failed to update project: %w", err)
 }

 slog.Info("Project updated successfully", "project_id", projectID, "user_id", userID)

 // Get owner details for response
 owner, err := s.userRepo.GetByID(ctx, project.OwnerID)
 if err != nil {
 slog.Warn("Failed to get owner details", "error", err, "owner_id", project.OwnerID)
 }

 resp := project.ToResponse()
 if owner != nil {
 ownerResp := owner.ToResponse()
 resp.Owner = &ownerResp
 }

 return &resp, nil
}

// DeleteProject deletes a project (owner only, cascade deletes tasks)
func (s *ProjectService) DeleteProject(ctx context.Context, projectID string, userID string) error {
 // Get existing project to check ownership
 project, err := s.projectRepo.GetByID(ctx, projectID)
 if err != nil {
 if errors.IsNotFound(err) {
 return fmt.Errorf("%w: project not found", errors.ErrNotFound)
 }
 slog.Error("Failed to get project for deletion", "error", err, "project_id", projectID)
 return fmt.Errorf("failed to get project: %w", err)
 }

 // Check ownership
 if project.OwnerID != userID {
 slog.Warn("Unauthorized project delete attempt", "project_id", projectID, "user_id", userID, "owner_id", project.OwnerID)
 return fmt.Errorf("%w: only the project owner can delete this project", errors.ErrUnauthorized)
 }

 // Delete project (tasks are cascade deleted via foreign key constraint)
 if err := s.projectRepo.Delete(ctx, projectID); err != nil {
 slog.Error("Failed to delete project", "error", err, "project_id", projectID)
 return fmt.Errorf("failed to delete project: %w", err)
 }

 slog.Info("Project deleted successfully", "project_id", projectID, "user_id", userID)
 return nil
}

// userHasAccess checks if a user has access to a project (owner or has assigned tasks)
func (s *ProjectService) userHasAccess(ctx context.Context, projectID, userID string) (bool, error) {
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

// validateCreateRequest validates a project creation request
func (s *ProjectService) validateCreateRequest(req model.ProjectCreateRequest) error {
 if req.Name == "" {
 return errors.NewValidationError("name", "name is required")
 }
 if len(req.Name) < 2 {
 return errors.NewValidationError("name", "name must be at least 2 characters")
 }
 if len(req.Name) > 255 {
 return errors.NewValidationError("name", "name must not exceed 255 characters")
 }
 return nil
}

// ProjectDetailResponse represents a project with its tasks
type ProjectDetailResponse struct {
 Project model.ProjectResponse `json:"project"`
 Tasks   []model.TaskResponse  `json:"tasks"`
}
