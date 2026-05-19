DROP TABLE IF EXISTS todo.tasks;
DROP EXTENSION IF EXISTS pg_trgm;
DROP INDEX IF EXISTS idx_tasks_title_trgm;
