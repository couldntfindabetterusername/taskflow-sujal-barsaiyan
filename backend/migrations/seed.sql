-- Seed Data Script for TaskFlow
-- This script is idempotent - can be run multiple times safely
--
-- Test User Credentials:
--   Email: test@example.com
--   Password: password123
--
-- To run manually:
--   docker compose exec db psql -U taskflow -d taskflow -f /docker-entrypoint-initdb.d/seed.sql
-- Or:
--   cat backend/migrations/seed.sql | docker compose exec -T db psql -U taskflow -d taskflow

-- Fixed UUIDs for repeatability
-- Test User: 11111111-1111-1111-1111-111111111111
-- Test Project: 22222222-2222-2222-2222-222222222222
-- Task 1 (todo): 33333333-3333-3333-3333-333333333331
-- Task 2 (in_progress): 33333333-3333-3333-3333-333333333332
-- Task 3 (done): 33333333-3333-3333-3333-333333333333

-- Insert test user (password: password123, bcrypt cost 12)
-- Hash generated with: golang.org/x/crypto/bcrypt.GenerateFromPassword([]byte("password123"), 12)
INSERT INTO users (id, name, email, password_hash, created_at)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'Test User',
    'test@example.com',
    '$2a$12$ETGBtzNfNUAd/.D8Vr9RmefmP1lPhKPMo9HHPlhjRmQVAybmuldtm',
    NOW()
)
ON CONFLICT (email) DO UPDATE SET
    name = EXCLUDED.name,
    password_hash = EXCLUDED.password_hash;

-- Insert test project owned by test user
INSERT INTO projects (id, name, description, owner_id, created_at)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    'Sample Project',
    'A sample project with test tasks for demonstration purposes.',
    '11111111-1111-1111-1111-111111111111',
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description;

-- Insert Task 1: Todo status, High priority, assigned to test user
INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
VALUES (
    '33333333-3333-3333-3333-333333333331',
    'Set up development environment',
    'Install all necessary tools and dependencies for the project.',
    'todo',
    'high',
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    NOW() + INTERVAL '7 days',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    priority = EXCLUDED.priority,
    assignee_id = EXCLUDED.assignee_id,
    due_date = EXCLUDED.due_date,
    updated_at = NOW();

-- Insert Task 2: In Progress status, Medium priority, assigned to test user
INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
VALUES (
    '33333333-3333-3333-3333-333333333332',
    'Implement user authentication',
    'Add login and registration functionality with JWT tokens.',
    'in_progress',
    'medium',
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    NOW() + INTERVAL '14 days',
    NOW() - INTERVAL '2 days',
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    priority = EXCLUDED.priority,
    assignee_id = EXCLUDED.assignee_id,
    due_date = EXCLUDED.due_date,
    updated_at = NOW();

-- Insert Task 3: Done status, Low priority, unassigned
INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
VALUES (
    '33333333-3333-3333-3333-333333333333',
    'Create project documentation',
    'Write README and API documentation for the project.',
    'done',
    'low',
    '22222222-2222-2222-2222-222222222222',
    NULL,
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '1 day'
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    priority = EXCLUDED.priority,
    assignee_id = EXCLUDED.assignee_id,
    due_date = EXCLUDED.due_date,
    updated_at = NOW();

-- Verify seed data
DO $$
DECLARE
    user_count INTEGER;
    project_count INTEGER;
    task_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO user_count FROM users WHERE email = 'test@example.com';
    SELECT COUNT(*) INTO project_count FROM projects WHERE id = '22222222-2222-2222-2222-222222222222';
    SELECT COUNT(*) INTO task_count FROM tasks WHERE project_id = '22222222-2222-2222-2222-222222222222';

    RAISE NOTICE 'Seed data verification:';
    RAISE NOTICE '  - Test user exists: %', user_count > 0;
    RAISE NOTICE '  - Test project exists: %', project_count > 0;
    RAISE NOTICE '  - Tasks in project: %', task_count;
END $$;
