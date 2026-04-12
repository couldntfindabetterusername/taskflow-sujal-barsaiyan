package handler

import (
 "encoding/json"
 "log/slog"
 "net/http"

 "github.com/taskflow/backend/internal/errors"
 "github.com/taskflow/backend/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
 authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
 return &AuthHandler{
 authService: authService,
 }
}

// ErrorResponse represents an error response
type ErrorResponse struct {
 Error string `json:"error"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
 var req service.RegisterRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 slog.Warn("Failed to decode register request", "error", err)
 respondWithError(w, http.StatusBadRequest, "invalid request body")
 return
 }

 resp, err := h.authService.Register(r.Context(), req)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusCreated, resp)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
 var req service.LoginRequest
 if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
 slog.Warn("Failed to decode login request", "error", err)
 respondWithError(w, http.StatusBadRequest, "invalid request body")
 return
 }

 resp, err := h.authService.Login(r.Context(), req)
 if err != nil {
 handleServiceError(w, err)
 return
 }

 respondWithJSON(w, http.StatusOK, resp)
}

// handleServiceError handles errors from the service layer and responds with appropriate HTTP status codes
func handleServiceError(w http.ResponseWriter, err error) {
 switch {
 case errors.IsValidation(err):
 respondWithError(w, http.StatusBadRequest, err.Error())
 case errors.IsDuplicate(err):
 respondWithError(w, http.StatusConflict, "email already registered")
 case errors.IsInvalidCredentials(err):
 respondWithError(w, http.StatusUnauthorized, "invalid email or password")
 case errors.IsNotFound(err):
 respondWithError(w, http.StatusNotFound, "resource not found")
 case errors.IsUnauthorized(err):
 respondWithError(w, http.StatusForbidden, "unauthorized")
 default:
 slog.Error("Internal server error", "error", err)
 respondWithError(w, http.StatusInternalServerError, "internal server error")
 }
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
 respondWithJSON(w, code, ErrorResponse{Error: message})
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(code)
 if err := json.NewEncoder(w).Encode(payload); err != nil {
 slog.Error("Failed to encode JSON response", "error", err)
 }
}
