// User types
export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

// Project types
export interface Project {
  id: string;
  name: string;
  description: string | null;
  owner_id: string;
  owner?: User;
  created_at: string;
}

export interface ProjectCreateRequest {
  name: string;
  description?: string;
}

export interface ProjectUpdateRequest {
  name?: string;
  description?: string;
}

export interface ProjectDetailResponse {
  project: Project;
  tasks: Task[];
}

// Task types
export type TaskStatus = 'todo' | 'in_progress' | 'done';
export type TaskPriority = 'low' | 'medium' | 'high';

export interface Task {
  id: string;
  title: string;
  description: string | null;
  status: TaskStatus;
  priority: TaskPriority;
  project_id: string;
  assignee_id: string | null;
  assignee?: User;
  due_date: string | null;
  created_at: string;
  updated_at: string;
}

export interface TaskCreateRequest {
  title: string;
  description?: string;
  status: TaskStatus;
  priority: TaskPriority;
  assignee_id?: string;
  due_date?: string;
}

export interface TaskUpdateRequest {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  assignee_id?: string;
  due_date?: string;
}

// API Error
export interface ApiError {
  error: string;
}
