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

// ProjectHandler handles project-related HTTP requests
type ProjectHandler struct {
 projectService *service.ProjectService
}

// NewProjectHandler creates a new ProjectHandler instance
func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
 return &ProjectHandler{
 projectService: projectService,
 }
}

// List handles GET /projects - list all accessible projects
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 // Parse pagination parameters
 limit := 50
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

 projects, err := h.projectService.ListProjects(r.Context(), userID, limit, offset)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 // Return empty array instead of null
 if projects == nil {
 projects = []model.ProjectResponse{}
 }

 respondWithJSON(w, http.StatusOK, projects)
}

// Create handles POST /projects - create a new project
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
 userID, ok := middleware.GetUserIDFromContext(r.Context())
 if !ok {
 respondWithError(w, http.StatusUnauthorized, "user not authenticated")
 return
 }

 var req model.ProjectCreateRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 slog.Warn("Failed to decode project create request", "error", err)
 respondWithError(w, http.StatusBadRequest, "invalid request body")
 return
 }

 project, err := h.projectService.CreateProject(r.Context(), req, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusCreated, project)
}

// Get handles GET /projects/:id - get project details with tasks
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
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

 project, err := h.projectService.GetProject(r.Context(), projectID, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusOK, project)
}

// Update handles PATCH /projects/:id - update a project
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
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

 var req model.ProjectUpdateRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 slog.Warn("Failed to decode project update request", "error", err)
 respondWithError(w, http.StatusBadRequest, "invalid request body")
 return
 }

 project, err := h.projectService.UpdateProject(r.Context(), projectID, req, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusOK, project)
}

// Delete handles DELETE /projects/:id - delete a project
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

 if err := h.projectService.DeleteProject(r.Context(), projectID, userID); err != nil {
 handleServiceError(w, err)
 return
 }

 w.WriteHeader(http.StatusNoContent)
}

// GetProjectStats gets statistics for a project
func (h *ProjectHandler) GetProjectStats(w http.ResponseWriter, r *http.Request) {
 userID := r.Context().Value("user_id").(string)

 projectID := chi.URLParam(r, "id")
 if projectID == "" {
 respondWithError(w, http.StatusBadRequest, "project ID is required")
 return
 }

 stats, err := h.projectService.GetProjectStats(r.Context(), projectID, userID)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusOK, stats)
}
