package service

import (
 "context"
 "fmt"
 "log/slog"
 "regexp"
 "time"

 "github.com/golang-jwt/jwt/v5"
 "github.com/taskflow/backend/internal/errors"
 "github.com/taskflow/backend/internal/model"
 "github.com/taskflow/backend/internal/repository"
 "golang.org/x/crypto/bcrypt"
)

const (
 // BcryptCost is the cost factor for bcrypt hashing
 BcryptCost = 12

 // JWTExpiryDuration is the expiry duration for JWT tokens (24 hours)
 JWTExpiryDuration = 24 * time.Hour
)

var (
 // emailRegex is a simple email validation regex
 emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
 UserID string `json:"user_id"`
 Email  string `json:"email"`
 jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
 userRepo  repository.UserRepository
 jwtSecret string
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
 return &AuthService{
 userRepo:  userRepo,
 jwtSecret: jwtSecret,
 }
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
 Name     string `json:"name"`
 Email    string `json:"email"`
 Password string `json:"password"`
}

// LoginRequest represents a login request
type LoginRequest struct {
 Email    string `json:"email"`
 Password string `json:"password"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
 Token string              `json:"token"`
 User  model.UserResponse  `json:"user"`
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
 // Validate input
 if err := s.validateRegisterRequest(req); err != nil {
 return nil, err
 }

 // Hash password
 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), BcryptCost)
 if err != nil {
 slog.Error("Failed to hash password", "error", err)
 return nil, fmt.Errorf("failed to hash password: %w", err)
 }

 // Create user
 user := &model.User{
 Name:         req.Name,
 Email:        req.Email,
 PasswordHash: string(hashedPassword),
 }

 if err := s.userRepo.Create(ctx, user); err != nil {
 if errors.IsDuplicate(err) {
 return nil, fmt.Errorf("%w: email already registered", errors.ErrDuplicate)
 }
 slog.Error("Failed to create user", "error", err)
 return nil, fmt.Errorf("failed to create user: %w", err)
 }

 slog.Info("User registered successfully", "user_id", user.ID, "email", user.Email)

 // Generate JWT token
 token, err := s.generateToken(user.ID, user.Email)
 if err != nil {
 return nil, err
 }

 return &AuthResponse{
 Token: token,
 User:  user.ToResponse(),
 }, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
 // Validate input
 if err := s.validateLoginRequest(req); err != nil {
 return nil, err
 }

 // Get user by email
 user, err := s.userRepo.GetByEmail(ctx, req.Email)
 if err != nil {
 if errors.IsNotFound(err) {
 return nil, fmt.Errorf("%w: invalid email or password", errors.ErrInvalidCredentials)
 }
 slog.Error("Failed to get user by email", "error", err, "email", req.Email)
 return nil, fmt.Errorf("failed to get user: %w", err)
 }

 // Compare password
 if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
 slog.Warn("Invalid password attempt", "user_id", user.ID, "email", user.Email)
 return nil, fmt.Errorf("%w: invalid email or password", errors.ErrInvalidCredentials)
 }

 slog.Info("User logged in successfully", "user_id", user.ID, "email", user.Email)

 // Generate JWT token
 token, err := s.generateToken(user.ID, user.Email)
 if err != nil {
 return nil, err
 }

 return &AuthResponse{
 Token: token,
 User:  user.ToResponse(),
 }, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(userID, email string) (string, error) {
 claims := JWTClaims{
 UserID: userID,
 Email:  email,
 RegisteredClaims: jwt.RegisteredClaims{
 ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpiryDuration)),
 IssuedAt:  jwt.NewNumericDate(time.Now()),
 },
 }

 token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
 tokenString, err := token.SignedString([]byte(s.jwtSecret))
 if err != nil {
 slog.Error("Failed to sign JWT token", "error", err)
 return "", fmt.Errorf("failed to generate token: %w", err)
 }

 return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
 token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
 // Validate signing method
 if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
 return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
 }
 return []byte(s.jwtSecret), nil
 })

 if err != nil {
 return nil, fmt.Errorf("%w: %v", errors.ErrInvalidToken, err)
 }

 claims, ok := token.Claims.(*JWTClaims)
 if !ok || !token.Valid {
 return nil, fmt.Errorf("%w: invalid token claims", errors.ErrInvalidToken)
 }

 return claims, nil
}

// validateRegisterRequest validates a registration request
func (s *AuthService) validateRegisterRequest(req RegisterRequest) error {
 if req.Name == "" {
 return errors.NewValidationError("name", "name is required")
 }

 if len(req.Name) < 2 {
 return errors.NewValidationError("name", "name must be at least 2 characters")
 }

 if req.Email == "" {
 return errors.NewValidationError("email", "email is required")
 }

 if !emailRegex.MatchString(req.Email) {
 return errors.NewValidationError("email", "invalid email format")
 }

 if req.Password == "" {
 return errors.NewValidationError("password", "password is required")
 }

 if len(req.Password) < 8 {
 return errors.NewValidationError("password", "password must be at least 8 characters")
 }

 return nil
}

// validateLoginRequest validates a login request
func (s *AuthService) validateLoginRequest(req LoginRequest) error {
 if req.Email == "" {
 return errors.NewValidationError("email", "email is required")
 }

 if !emailRegex.MatchString(req.Email) {
 return errors.NewValidationError("email", "invalid email format")
 }

 if req.Password == "" {
 return errors.NewValidationError("password", "password is required")
 }

 return nil
}
