CREATE TABLE IF NOT EXISTS todo.tasks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    description TEXT,
    category TEXT,
    tags TEXT[],
    created_at timestamptz DEFAULT NOW()
)