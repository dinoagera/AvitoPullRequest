-- -- +goose Up
CREATE TABLE IF NOT EXISTS pull_requests (
    pr_id TEXT PRIMARY KEY,              
    title TEXT NOT NULL,                 
    author_id TEXT NOT NULL,             
    status TEXT NOT NULL DEFAULT 'OPEN', 
    assigned_reviewers TEXT[],          
    created_at TIMESTAMP WITH TIME ZONE, 
    merged_at TIMESTAMP WITH TIME ZONE  
);
ALTER TABLE pull_requests
ADD CONSTRAINT fk_author_id
FOREIGN KEY (author_id) REFERENCES users(user_id);
CREATE INDEX idx_assigned_reviewers ON pull_requests USING GIN (assigned_reviewers);
-- -- +goose Down
DROP INDEX IF EXISTS idx_assigned_reviewers;
DROP TABLE IF EXISTS pull_requests;
