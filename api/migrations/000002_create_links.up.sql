CREATE TABLE IF NOT EXISTS links (
    id           SERIAL PRIMARY KEY,
    project_id   INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    original_url TEXT NOT NULL,
    title        VARCHAR(255) NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    og_image     TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
