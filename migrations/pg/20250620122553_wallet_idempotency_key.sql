-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS idempotency_key TEXT UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
