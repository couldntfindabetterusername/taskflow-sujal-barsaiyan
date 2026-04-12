package model

import (
 "fmt"
 "time"
)

// TaskStatus represents the possible status values for a task
type TaskStatus string

const (
 TaskStatusTodo       TaskStatus = "todo"
 TaskStatusInProgress TaskStatus = "in_progress"
 TaskStatusDone       TaskStatus = "done"
)

// Valid returns true if the task status is valid
func (s TaskStatus) Valid() bool {
 switch s {
 case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
 return true
 }
 return false
}

// String returns the string representation of the task status
func (s TaskStatus) String() string {
 return string(s)
}

// TaskPriority represents the possible priority values for a task
type TaskPriority string

const (
 TaskPriorityLow    TaskPriority = "low"
 TaskPriorityMedium TaskPriority = "medium"
 TaskPriorityHigh   TaskPriority = "high"
)

// Valid returns true if the task priority is valid
func (p TaskPriority) Valid() bool {
 switch p {
 case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh:
 return true
 }
 return false
}

// String returns the string representation of the task priority
func (p TaskPriority) String() string {
 return string(p)
}

// Task represents a task in the system
type Task struct {
 ID          string        `json:"id"`
 Title       string        `json:"title"`
 Description *string       `json:"description"`
 Status      TaskStatus    `json:"status"`
 Priority    TaskPriority  `json:"priority"`
 ProjectID   string        `json:"project_id"`
 AssigneeID  *string       `json:"assignee_id"`
 DueDate     *time.Time    `json:"due_date"`
 CreatedAt   time.Time     `json:"created_at"`
 UpdatedAt   time.Time     `json:"updated_at"`
}

// TaskCreateRequest represents the request body for creating a new task
type TaskCreateRequest struct {
 Title       string       `json:"title"`
 Description *string      `json:"description,omitempty"`
 Status      TaskStatus   `json:"status"`
 Priority    TaskPriority `json:"priority"`
 AssigneeID  *string      `json:"assignee_id,omitempty"`
 DueDate     *time.Time   `json:"due_date,omitempty"`
}

// Validate validates the task create request
func (r *TaskCreateRequest) Validate() error {
 if r.Title == "" {
 return fmt.Errorf("title is required")
 }
 if !r.Status.Valid() {
 return fmt.Errorf("invalid status: must be one of 'todo', 'in_progress', 'done'")
 }
 if !r.Priority.Valid() {
 return fmt.Errorf("invalid priority: must be one of 'low', 'medium', 'high'")
 }
 return nil
}

// TaskUpdateRequest represents the request body for updating a task
type TaskUpdateRequest struct {
 Title       *string       `json:"title,omitempty"`
 Description *string       `json:"description,omitempty"`
 Status      *TaskStatus   `json:"status,omitempty"`
 Priority    *TaskPriority `json:"priority,omitempty"`
 AssigneeID  *string       `json:"assignee_id,omitempty"`
 DueDate     *time.Time    `json:"due_date,omitempty"`
}

// Validate validates the task update request
func (r *TaskUpdateRequest) Validate() error {
 if r.Status != nil && !r.Status.Valid() {
 return fmt.Errorf("invalid status: must be one of 'todo', 'in_progress', 'done'")
 }
 if r.Priority != nil && !r.Priority.Valid() {
 return fmt.Errorf("invalid priority: must be one of 'low', 'medium', 'high'")
 }
 return nil
}

// TaskResponse represents a task in API responses with related entities
type TaskResponse struct {
 ID          string        `json:"id"`
 Title       string        `json:"title"`
 Description *string       `json:"description"`
 Status      TaskStatus    `json:"status"`
 Priority    TaskPriority  `json:"priority"`
 ProjectID   string        `json:"project_id"`
 AssigneeID  *string       `json:"assignee_id"`
 Assignee    *UserResponse `json:"assignee,omitempty"`
 DueDate     *time.Time    `json:"due_date"`
 CreatedAt   time.Time     `json:"created_at"`
 UpdatedAt   time.Time     `json:"updated_at"`
}

// ToResponse converts a Task to TaskResponse
func (t *Task) ToResponse() TaskResponse {
 return TaskResponse{
 ID:          t.ID,
 Title:       t.Title,
 Description: t.Description,
 Status:      t.Status,
 Priority:    t.Priority,
 ProjectID:   t.ProjectID,
 AssigneeID:  t.AssigneeID,
 DueDate:     t.DueDate,
 CreatedAt:   t.CreatedAt,
 UpdatedAt:   t.UpdatedAt,
 }
}

// TaskFilters represents filters for querying tasks
type TaskFilters struct {
 Status     *TaskStatus   `json:"status,omitempty"`
 Priority   *TaskPriority `json:"priority,omitempty"`
 AssigneeID *string       `json:"assignee_id,omitempty"`
}
