import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  Project,
  ProjectCreateRequest,
  ProjectUpdateRequest,
  ProjectDetailResponse,
  Task,
  TaskCreateRequest,
  TaskUpdateRequest,
  ApiError,
} from '../types';

// Create axios instance
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiError>) => {
    if (error.response?.status === 401) {
      // Clear token and redirect to login on unauthorized
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      // Only redirect if not already on auth pages
      if (!window.location.pathname.includes('/login') && !window.location.pathname.includes('/register')) {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authApi = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/login', data);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/auth/register', data);
    return response.data;
  },
};

// Projects API
export const projectsApi = {
  list: async (limit = 50, offset = 0): Promise<Project[]> => {
    const response = await api.get<Project[]>('/projects', {
      params: { limit, offset },
    });
    return response.data;
  },

  get: async (id: string): Promise<ProjectDetailResponse> => {
    const response = await api.get<ProjectDetailResponse>(`/projects/${id}`);
    return response.data;
  },

  create: async (data: ProjectCreateRequest): Promise<Project> => {
    const response = await api.post<Project>('/projects', data);
    return response.data;
  },

  update: async (id: string, data: ProjectUpdateRequest): Promise<Project> => {
    const response = await api.patch<Project>(`/projects/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/projects/${id}`);
  },
};

// Tasks API
export const tasksApi = {
  list: async (
    projectId: string,
    filters?: { status?: string; assignee?: string; priority?: string },
    limit = 100,
    offset = 0
  ): Promise<Task[]> => {
    const response = await api.get<Task[]>(`/projects/${projectId}/tasks`, {
      params: { ...filters, limit, offset },
    });
    return response.data;
  },

  get: async (id: string): Promise<Task> => {
    const response = await api.get<Task>(`/tasks/${id}`);
    return response.data;
  },

  create: async (projectId: string, data: TaskCreateRequest): Promise<Task> => {
    const response = await api.post<Task>(`/projects/${projectId}/tasks`, data);
    return response.data;
  },

  update: async (id: string, data: TaskUpdateRequest): Promise<Task> => {
    const response = await api.patch<Task>(`/tasks/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/tasks/${id}`);
  },
};

// Helper to extract error message
export const getErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<ApiError>;
    return axiosError.response?.data?.error || axiosError.message || 'An error occurred';
  }
  if (error instanceof Error) {
    return error.message;
  }
  return 'An unexpected error occurred';
};

export default api;
