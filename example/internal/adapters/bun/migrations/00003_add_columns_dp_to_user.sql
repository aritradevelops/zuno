-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN dp VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN dp;
-- +goose StatementEnd
