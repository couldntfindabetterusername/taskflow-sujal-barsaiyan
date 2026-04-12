package handler

import (
 "encoding/json"
 "log/slog"
 "net/http"
 "strconv"

 "github.com/go-chi/chi/v5"
 "github.com/taskflow/backend/internal/middleware"
 "github.com/taskflow/backend/internal/model"
 "github.com/taskflow/backend/internal/service"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
 taskService *service.TaskService
}

// NewTaskHandler creates a new TaskHandler instance
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
 return &TaskHandler{
 taskService: taskService,
 }
}

// List handles GET /projects/:id/tasks - list tasks in a project with optional filters
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 projectID := chi.URLParam(r, "id")
 if projectID == "" {
 respondWithError(w, http.StatusBadRequest, "project ID is required")
 return
 }

 // Parse filters from query parameters
 filters := &model.TaskFilters{}

 if statusStr := r.URL.Query().Get("status"); statusStr != "" {
 status := model.TaskStatus(statusStr)
 if !status.Valid() {
 respondWithError(w, http.StatusBadRequest, "invalid status: must be one of 'todo', 'in_progress', 'done'")
 return
 }
 filters.Status = &status
 }

 if assigneeID := r.URL.Query().Get("assignee"); assigneeID != "" {
 filters.AssigneeID = &assigneeID
 }

 if priorityStr := r.URL.Query().Get("priority"); priorityStr != "" {
 priority := model.TaskPriority(priorityStr)
 if !priority.Valid() {
 respondWithError(w, http.StatusBadRequest, "invalid priority: must be one of 'low', 'medium', 'high'")
 return
 }
 filters.Priority = &priority
 }

 // Parse pagination parameters
 limit := 100
 offset := 0

 if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
 if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
 limit = l
 }
 }

 if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
 if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
 offset = o
 }
 }

 tasks, err := h.taskService.ListTasks(r.Context(), projectID, userID, filters, limit, offset)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 // Return empty array instead of null
 if tasks == nil {
 tasks = []model.TaskResponse{}
 }

 respondWithJSON(w, http.StatusOK, tasks)
}

// Create handles POST /projects/:id/tasks - create a new task in a project
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 projectID := chi.URLParam(r, "id")
 if projectID == "" {
 respondWithError(w, http.StatusBadRequest, "project ID is required")
 return
 }

 var req model.TaskCreateRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 slog.Warn("Failed to decode task create request", "error", err)
 respondWithError(w, http.StatusBadRequest, "invalid request body")
 return
 }

 task, err := h.taskService.CreateTask(r.Context(), projectID, req, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusCreated, task)
}

// Get handles GET /tasks/:id - get a single task
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 taskID := chi.URLParam(r, "id")
 if taskID == "" {
 respondWithError(w, http.StatusBadRequest, "task ID is required")
 return
 }

 task, err := h.taskService.GetTask(r.Context(), taskID, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusOK, task)
}

// Update handles PATCH /tasks/:id - update a task
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 taskID := chi.URLParam(r, "id")
 if taskID == "" {
 respondWithError(w, http.StatusBadRequest, "task ID is required")
 return
 }

 var req model.TaskUpdateRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 slog.Warn("Failed to decode task update request", "error", err)
 respondWithError(w, http.StatusBadRequest, "invalid request body")
 return
 }

 task, err := h.taskService.UpdateTask(r.Context(), taskID, req, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusOK, task)
}

// Delete handles DELETE /tasks/:id - delete a task
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 taskID := chi.URLParam(r, "id")
 if taskID == "" {
 respondWithError(w, http.StatusBadRequest, "task ID is required")
 return
 }

 if err := h.taskService.DeleteTask(r.Context(), taskID, userID); err != nil {
 handleServiceError(w, err)
 return
 }

 w.WriteHeader(http.StatusNoContent)
}
