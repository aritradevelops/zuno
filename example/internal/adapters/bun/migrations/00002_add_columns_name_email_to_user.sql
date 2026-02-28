-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN name;
ALTER TABLE users DROP COLUMN email;
-- +goose StatementEnd
