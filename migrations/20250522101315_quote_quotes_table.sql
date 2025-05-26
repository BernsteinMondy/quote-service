-- +goose Up
-- +goose StatementBegin
CREATE TABLE quote.quotes
(
    id     uuid PRIMARY KEY,
    author text NOT NULL,
    quote  text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quote.quotes;
-- +goose StatementEnd
