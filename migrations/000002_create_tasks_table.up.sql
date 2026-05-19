CREATE TABLE IF NOT EXISTS todo.tasks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    description TEXT,
    category TEXT,
    tags TEXT[],
    created_at timestamptz DEFAULT NOW()
);
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_tasks_title_trgm ON todo.tasks USING gin(title gin_trgm_ops);
