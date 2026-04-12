-- Drop tasks table, indexes, trigger and function
DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_tasks_assignee_id;
DROP INDEX IF EXISTS idx_tasks_project_id;
DROP TABLE IF EXISTS tasks;
