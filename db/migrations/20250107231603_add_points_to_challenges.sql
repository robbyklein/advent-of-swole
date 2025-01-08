-- +goose Up
-- +goose StatementBegin
ALTER TABLE challenges ADD COLUMN points INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE challenges DROP COLUMN points;
-- +goose StatementEnd