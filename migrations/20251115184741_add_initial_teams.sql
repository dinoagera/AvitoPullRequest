-- -- +goose Up
CREATE TABLE IF NOT EXISTS teams(
    name TEXT PRIMARY KEY        
);
-- -- +goose Down
DROP TABLE IF EXISTS teams;
