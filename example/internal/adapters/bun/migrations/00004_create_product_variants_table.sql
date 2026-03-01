-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_variants (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	uid UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by UUID NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	updated_by UUID NOT NULL,
	deleted_at TIMESTAMPTZ,
	deleted_by UUID
);
CREATE INDEX idx_product_variants_uid ON product_variants(uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_variants_uid;
DROP TABLE IF EXISTS product_variants;
-- +goose StatementEnd