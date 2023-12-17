-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files;
-- +goose StatementEnd
