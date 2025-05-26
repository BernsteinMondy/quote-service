-- +goose Up
-- +goose StatementBegin
CREATE INDEX index_quote_quotes_author ON quote.quotes(author);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX quote.index_quote_quotes_author;
-- +goose StatementEnd
