package middleware

import (
 "context"
 "log/slog"
 "net/http"
 "strings"

 "github.com/taskflow/backend/internal/service"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
 // UserIDKey is the context key for the user ID
 UserIDKey contextKey = "user_id"

 // UserEmailKey is the context key for the user email
 UserEmailKey contextKey = "user_email"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
 return func(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 // Extract token from Authorization header
 authHeader := r.Header.Get("Authorization")
 if authHeader == "" {
 respondUnauthorized(w, "missing authorization header")
 return
 }

 // Check if the header starts with "Bearer "
 parts := strings.Split(authHeader, " ")
 if len(parts) != 2 || parts[0] != "Bearer" {
 respondUnauthorized(w, "invalid authorization header format")
 return
 }

 tokenString := parts[1]

 // Validate token
 claims, err := authService.ValidateToken(tokenString)
 if err != nil {
 slog.Warn("Invalid token", "error", err)
 respondUnauthorized(w, "invalid or expired token")
 return
 }

 // Add user info to context
 ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
 ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

 // Call next handler with updated context
 next.ServeHTTP(w, r.WithContext(ctx))
 })
 }
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
 userID, ok := ctx.Value(UserIDKey).(string)
 return userID, ok
}

// GetUserEmailFromContext extracts the user email from the request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
 email, ok := ctx.Value(UserEmailKey).(string)
 return email, ok
}

// respondUnauthorized sends an unauthorized response
func respondUnauthorized(w http.ResponseWriter, message string) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(http.StatusUnauthorized)
 w.Write([]byte(`{"error":"` + message + `"}`))
}
