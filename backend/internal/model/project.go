package model

import "time"

// Project represents a project in the system
type Project struct {
 ID          string    `json:"id"`
 Name        string    `json:"name"`
 Description *string   `json:"description"`
 OwnerID     string    `json:"owner_id"`
 CreatedAt   time.Time `json:"created_at"`
}

// ProjectCreateRequest represents the request body for creating a new project
type ProjectCreateRequest struct {
 Name        string  `json:"name"`
 Description *string `json:"description,omitempty"`
}

// ProjectUpdateRequest represents the request body for updating a project
type ProjectUpdateRequest struct {
 Name        *string `json:"name,omitempty"`
 Description *string `json:"description,omitempty"`
}

// ProjectResponse represents a project in API responses with owner details
type ProjectResponse struct {
 ID          string        `json:"id"`
 Name        string        `json:"name"`
 Description *string       `json:"description"`
 OwnerID     string        `json:"owner_id"`
 Owner       *UserResponse `json:"owner,omitempty"`
 CreatedAt   time.Time     `json:"created_at"`
}

// ToResponse converts a Project to ProjectResponse
func (p *Project) ToResponse() ProjectResponse {
 return ProjectResponse{
 ID:          p.ID,
 Name:        p.Name,
 Description: p.Description,
 OwnerID:     p.OwnerID,
 CreatedAt:   p.CreatedAt,
 }
}
