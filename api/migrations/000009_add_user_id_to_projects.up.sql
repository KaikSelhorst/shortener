ALTER TABLE projects ADD COLUMN user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_projects_user_id ON projects(user_id);
