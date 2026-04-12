package main

import (
 "context"
 "fmt"
 "log/slog"
 "net/http"
 "os"
 "os/signal"
 "syscall"
 "time"

 "github.com/go-chi/chi/v5"
 "github.com/go-chi/chi/v5/middleware"
 "github.com/jackc/pgx/v5/pgxpool"
 "github.com/taskflow/backend/internal/config"
 "github.com/taskflow/backend/internal/handler"
 custommiddleware "github.com/taskflow/backend/internal/middleware"
 "github.com/taskflow/backend/internal/repository/postgres"
 "github.com/taskflow/backend/internal/service"
)

func main() {
 logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
 slog.SetDefault(logger)

 cfg, err := config.Load()
 if err != nil {
 slog.Error("Failed to load configuration", "error", err)
 os.Exit(1)
 }

 db, err := connectWithRetry(cfg.DatabaseURL, 5, 2*time.Second)
 if err != nil {
 slog.Error("Failed to connect to database", "error", err)
 os.Exit(1)
 }
 defer db.Close()

 slog.Info("Successfully connected to database")

 // Initialize repositories
 userRepo := postgres.NewUserRepository(db)
 projectRepo := postgres.NewProjectRepository(db)
 taskRepo := postgres.NewTaskRepository(db)

 // Initialize services
 authService := service.NewAuthService(userRepo, cfg.JWTSecret)
 projectService := service.NewProjectService(projectRepo, taskRepo, userRepo)
 taskService := service.NewTaskService(taskRepo, projectRepo, userRepo)

 // Initialize handlers
 authHandler := handler.NewAuthHandler(authService)
 projectHandler := handler.NewProjectHandler(projectService)
 taskHandler := handler.NewTaskHandler(taskService)

 // Setup router
 r := chi.NewRouter()

 r.Use(middleware.Logger)
 r.Use(middleware.Recoverer)
 r.Use(middleware.RequestID)
 r.Use(middleware.RealIP)
 r.Use(corsMiddleware)

 // Public routes
 r.Get("/health", healthHandler)

 // Auth routes (public)
 r.Route("/auth", func(r chi.Router) {
 r.Post("/register", authHandler.Register)
 r.Post("/login", authHandler.Login)
 })

 // Protected routes
 r.Group(func(r chi.Router) {
 r.Use(custommiddleware.AuthMiddleware(authService))

 // Projects routes
 r.Route("/projects", func(r chi.Router) {
 r.Get("/", projectHandler.List)
 r.Post("/", projectHandler.Create)
 r.Get("/{id}", projectHandler.Get)
 r.Patch("/{id}", projectHandler.Update)
 r.Delete("/{id}", projectHandler.Delete)

 // Tasks routes nested under projects
 r.Get("/{id}/tasks", taskHandler.List)
 r.Post("/{id}/tasks", taskHandler.Create)
 })

 // Tasks routes (standalone for update/delete)
 r.Route("/tasks", func(r chi.Router) {
 r.Get("/{id}", taskHandler.Get)
 r.Patch("/{id}", taskHandler.Update)
 r.Delete("/{id}", taskHandler.Delete)
 })
 })

 addr := fmt.Sprintf(":%s", cfg.ServerPort)
 server := &http.Server{
 Addr:         addr,
 Handler:      r,
 ReadTimeout:  15 * time.Second,
 WriteTimeout: 15 * time.Second,
 IdleTimeout:  60 * time.Second,
 }

 go func() {
 slog.Info("Starting server", "address", addr)
 if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
 slog.Error("Server failed to start", "error", err)
 os.Exit(1)
 }
 }()

 quit := make(chan os.Signal, 1)
 signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
 <-quit

 slog.Info("Shutting down server...")

 ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
 defer cancel()

 if err := server.Shutdown(ctx); err != nil {
 slog.Error("Server forced to shutdown", "error", err)
 os.Exit(1)
 }

 slog.Info("Server exited gracefully")
}

func connectWithRetry(databaseURL string, maxAttempts int, delay time.Duration) (*pgxpool.Pool, error) {
 var db *pgxpool.Pool
 var err error

 for i := 1; i <= maxAttempts; i++ {
 slog.Info("Attempting to connect to database", "attempt", i, "max_attempts", maxAttempts)

 ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
 db, err = pgxpool.New(ctx, databaseURL)
 cancel()

 if err == nil {
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 err = db.Ping(ctx)
 cancel()

 if err == nil {
 return db, nil
 }
 }

 if i < maxAttempts {
 slog.Warn("Failed to connect to database, retrying...", "error", err, "retry_in", delay)
 time.Sleep(delay)
 }
 }

 return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}

func corsMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Access-Control-Allow-Origin", "*")
 w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
 w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
 w.Header().Set("Access-Control-Max-Age", "3600")

 if r.Method == "OPTIONS" {
 w.WriteHeader(http.StatusOK)
 return
 }

 next.ServeHTTP(w, r)
 })
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json")
 w.WriteHeader(http.StatusOK)
 w.Write([]byte(`{"status":"ok"}`))
}
