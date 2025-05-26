-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA quote;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA quote;
-- +goose StatementEnd
