-- +goose Up
CREATE TABLE translations (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    api_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE translations;