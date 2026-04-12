# TaskFlow - Full-Stack Task Management Application

A modern, production-ready task management system demonstrating best practices in full-stack development with a focus on clean architecture, security, and developer experience.

## 📋 Overview

TaskFlow is a task management application that enables teams to organize work efficiently. Users can create projects, manage tasks with different statuses and priorities, and track progress—all through an intuitive, responsive web interface.

### Tech Stack

**Backend:**
- Go 1.22 with Chi router for lightweight, fast HTTP routing
- PostgreSQL 15 with pgx/v5 driver for robust data persistence
- golang-migrate for version-controlled database migrations
- JWT (golang-jwt/jwt/v5) for stateless authentication
- bcrypt for secure password hashing (cost factor: 12)
- Structured logging with Go's slog package

**Frontend:**
- React 18 with TypeScript for type-safe component development
- Vite for blazing-fast development and optimized production builds
- Material UI (MUI) v5 for consistent, accessible component design
- TanStack Query (React Query) for server state management and caching
- Axios for HTTP client with request/response interceptors
- React Router v6 for declarative routing

**Infrastructure:**
- Docker with multi-stage builds for minimal production images
- Docker Compose for orchestrating services (PostgreSQL, backend, frontend)
- Nginx as a static file server for the frontend in production
- Health checks and dependency management between services

---

## 🏗️ Architecture Decisions

### Why Material UI?

MUI provides a comprehensive set of accessible, well-designed components out of the box. This accelerates development while maintaining professional design standards and accessibility.

**Tradeoffs:**
- ✅ Consistent design system with minimal custom CSS
- ✅ Accessibility (ARIA) built-in
- ✅ Responsive components with built-in breakpoint system
- ⚠️ Larger bundle size than minimal UI libraries
- ⚠️ Customization can be complex for highly branded designs

### Why TanStack Query?

TanStack Query (React Query) handles server state management elegantly, providing automatic caching, background refetching, and optimistic updates. This eliminates boilerplate and improves user experience.

**Tradeoffs:**
- ✅ Automatic request deduplication and caching
- ✅ Optimistic updates with automatic rollback
- ✅ Background refetching keeps data fresh
- ⚠️ Additional library dependency

### Repository Pattern

The backend uses a three-layer architecture (handler → service → repository) to separate concerns:
- **Handlers**: HTTP request/response handling
- **Services**: Business logic and authorization
- **Repositories**: Data access and persistence

**Tradeoffs:**
- ✅ Easy to test each layer in isolation
- ✅ Clean separation of concerns
- ✅ Easy to swap implementations (e.g., mock repository for tests)
- ⚠️ More files and interfaces to maintain
- ⚠️ Can feel over-engineered for very simple CRUD operations

---

## 🚀 Running Locally

### Prerequisites

- **Docker** and **Docker Compose** (Docker Desktop includes both)
- **Git** (to clone the repository)

No need to install Go, Node.js, PostgreSQL, or any other dependencies—Docker handles everything!

### Quick Start

1. **Clone the repository:**

```bash
git clone <repository-url>
cd "taskflow-sujal-barsaiyan"
```

2. **Create environment file:**

```bash
# Copy the example environment file
cp .env.example .env
```

The `.env.example` file contains secure defaults for local development. No changes needed for basic usage.

3. **Start the application:**

```bash
docker compose up --build
```

This single command will:
- Build the backend Go application
- Build the frontend React application
- Start PostgreSQL database
- Run database migrations automatically
- Seed the database with test data
- Start all services with proper networking

4. **Access the application:**

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Backend Health Check**: http://localhost:8080/health

5. **Login with test credentials:**

```
Email: test@example.com
Password: password123
```

### Stopping the Application

```bash
docker compose down
```

### Viewing Logs

```bash
docker compose logs -f
```

---

## 🗄️ Database Migrations

Migrations run **automatically** when the backend container starts. No manual intervention required!

### Migration Files

Located in `backend/migrations/`:
- `000001_create_users_table.up.sql` / `.down.sql`
- `000002_create_projects_table.up.sql` / `.down.sql`
- `000003_create_tasks_table.up.sql` / `.down.sql`
- `seed.sql` (test data)

### Manual Migration Commands (if needed)

Run migrations manually (only if needed):

```bash
# Run all up migrations
docker compose exec backend migrate -path /app/migrations -database "${DATABASE_URL}" up

# Rollback last migration
docker compose exec backend migrate -path /app/migrations -database "${DATABASE_URL}" down 1

# Check migration version
docker compose exec backend migrate -path /app/migrations -database "${DATABASE_URL}" version
```

### Seed Test Data

Test data is automatically loaded, but you can re-run it:

```bash
cat backend/migrations/seed.sql | docker compose exec -T db psql -U taskflow -d taskflow
```

**Test Data Includes:**
- 1 test user (test@example.com / password123)
- 1 project owned by the test user
- 3 tasks with different statuses and priorities

---

## 📡 API Reference

### Base URL

```
http://localhost:8080
```

### Authentication Endpoints

#### Register User

```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-here",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-here",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Project Endpoints

**All project endpoints require authentication.** Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-token-here>
```

#### List Projects

```http
GET /projects
```

**Response (200 OK):**
```json
[
  {
    "id": "uuid-here",
    "name": "Website Redesign",
    "description": "Complete overhaul of company website",
    "owner_id": "owner-uuid",
    "owner_name": "John Doe",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

#### Create Project

```http
POST /projects
Content-Type: application/json

{
  "name": "Mobile App",
  "description": "iOS and Android app development"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid-here",
  "name": "Mobile App",
  "description": "iOS and Android app development",
  "owner_id": "your-uuid",
  "owner_name": "John Doe",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Get Project Details

```http
GET /projects/:id
```

**Response (200 OK):**
```json
{
  "id": "uuid-here",
  "name": "Website Redesign",
  "description": "Complete overhaul of company website",
  "owner_id": "owner-uuid",
  "owner_name": "John Doe",
  "created_at": "2024-01-15T10:30:00Z",
  "tasks": [
    {
      "id": "task-uuid",
      "title": "Design homepage mockup",
      "description": "Create Figma designs",
      "status": "in_progress",
      "priority": "high",
      "project_id": "uuid-here",
      "assignee_id": "user-uuid",
      "assignee_name": "Jane Smith",
      "due_date": "2024-02-01T00:00:00Z",
      "created_at": "2024-01-16T09:00:00Z",
      "updated_at": "2024-01-16T09:00:00Z"
    }
  ]
}
```

#### Update Project

```http
PATCH /projects/:id
Content-Type: application/json

{
  "name": "Website Redesign v2",
  "description": "Updated description"
}
```

**Response (200 OK):** Updated project object

#### Delete Project

```http
DELETE /projects/:id
```

**Response (204 No Content)**

**Note:** Deletes all tasks in the project (cascade delete).

#### Get Project Statistics ⭐

```http
GET /projects/:id/stats
```

**Response (200 OK):**
```json
{
  "project_id": "uuid-here",
  "total_tasks": 10,
  "completion_percentage": 30.0,
  "by_status": {
    "todo": 5,
    "in_progress": 2,
    "done": 3
  },
  "by_assignee": {
    "user-uuid-1": 4,
    "user-uuid-2": 3,
    "unassigned": 3
  }
}
```

**Authorization:** Only accessible to project members (owner or users with assigned tasks).

**Use Cases:**
- Dashboard displays
- Progress tracking
- Resource allocation insights
- Completion metrics

### Task Endpoints

#### List Tasks in Project

```http
GET /projects/:id/tasks?status=in_progress&assignee=me&priority=high
```

**Query Parameters:**
- `status` (optional): Filter by status (`todo`, `in_progress`, `done`)
- `assignee` (optional): Filter by assignee (`me` for current user, `none` for unassigned, or user UUID)
- `priority` (optional): Filter by priority (`low`, `medium`, `high`)
- `limit` (optional): Number of results (default: 100)
- `offset` (optional): Pagination offset (default: 0)

**Response (200 OK):**
```json
[
  {
    "id": "task-uuid",
    "title": "Design homepage mockup",
    "description": "Create Figma designs for new homepage",
    "status": "in_progress",
    "priority": "high",
    "project_id": "project-uuid",
    "assignee_id": "user-uuid",
    "assignee_name": "Jane Smith",
    "due_date": "2024-02-01T00:00:00Z",
    "created_at": "2024-01-16T09:00:00Z",
    "updated_at": "2024-01-16T14:30:00Z"
  }
]
```

#### Create Task

```http
POST /projects/:id/tasks
Content-Type: application/json

{
  "title": "Implement login page",
  "description": "Create React component for user login",
  "status": "todo",
  "priority": "high",
  "assignee_id": "user-uuid",
  "due_date": "2024-02-15T00:00:00Z"
}
```

**Response (201 Created):** Task object

#### Get Task Details

```http
GET /tasks/:id
```

**Response (200 OK):** Task object

#### Update Task

```http
PATCH /tasks/:id
Content-Type: application/json

{
  "status": "in_progress",
  "priority": "medium"
}
```

**Response (200 OK):** Updated task object

**Note:** Any field can be updated. Omitted fields remain unchanged.

#### Delete Task

```http
DELETE /tasks/:id
```

**Response (204 No Content)**

**Authorization:** Only the project owner can delete tasks.

### Error Responses

All errors follow this format:

```json
{
  "error": "Human-readable error message"
}
```

**Common Status Codes:**
- `400 Bad Request` - Invalid input or validation error
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource doesn't exist
- `409 Conflict` - Duplicate resource (e.g., email already exists)
- `500 Internal Server Error` - Unexpected server error

---

## 📱UI Features

### Responsive Design

The application is fully responsive and tested at the following breakpoints:

- **Mobile**: 375px (iPhone SE, Galaxy S8+)
- **Tablet**: 768px (iPad)
- **Desktop**: 1280px and above

### Dark Mode

The application includes a fully functional dark mode toggle. Click the sun/moon icon in the navigation bar to toggle between light and dark modes!

---

## 🎯 What I'd Do With More Time

1. **Comprehensive Testing**
   - Unit tests for Go services and repositories (using testify)
   - Frontend component tests with React Testing Library
   - E2E tests with Playwright for critical user flows
   - API integration tests
   - Target: 80%+ code coverage

2. **Real-time Collaboration**
   - WebSocket support for live task updates
   - See when team members are viewing/editing tasks
   - Real-time notifications for task assignments and status changes

3. **Advanced Task Features**
   - Task assignee selection in UI (currently tasks can only be assigned via API)
   - Subtasks and task dependencies
   - Task comments and activity history
   - File attachments (with S3 or similar storage)
   - Task templates for common workflows
   - Bulk task operations (select multiple, update status)

4. **Enhanced Security**
   - Rate limiting on API endpoints
   - CSRF protection
   - Email verification on registration
   - Password reset flow via email
   - Two-factor authentication (2FA)
   - Session management (token refresh, logout all devices)

5. **Performance Optimizations**
   - Database query optimization with EXPLAIN ANALYZE
   - Redis caching layer for frequently accessed data
   - Pagination for large task lists (currently loads all)
   - Lazy loading for task descriptions and details
   - Frontend code splitting for faster initial load

6. **User Experience**
   - Drag-and-drop task reordering (react-beautiful-dnd)
   - Kanban board view (in addition to list view)
   - Keyboard shortcuts (j/k navigation, ? for help)
   - Undo/redo for task operations
   - Export projects/tasks to CSV or JSON

7. **Team Collaboration**
   - Project team member management (invite users to projects)
   - Role-based access control (admin, member, viewer)
   - @mentions in task descriptions
   - Email notifications for task updates

8. **Analytics & Reporting**
   - Dashboard with project/task statistics
   - Burndown charts for sprint tracking
   - Velocity tracking over time
   - Task completion rates by user
   - Time tracking per task
   - Due date compliance metrics
   
9. **Additional Features**
    - Search functionality (full-text search on tasks/projects)
    - Custom fields per project
    - Task recurrence (daily, weekly, monthly)
    - Calendar view for due dates
    - Project archives (soft delete)
    - Audit log for compliance

