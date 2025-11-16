-- -- +goose Up
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);
-- -- +goose Down
DROP TABLE IF EXISTS users;
