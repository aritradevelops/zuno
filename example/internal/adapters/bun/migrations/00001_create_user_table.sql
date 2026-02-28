-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	uid UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by UUID NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	updated_by UUID NOT NULL,
	deleted_at TIMESTAMPTZ,
	deleted_by UUID
);
CREATE INDEX idx_users_uid ON users(uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_uid;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
