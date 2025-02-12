-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  codename VARCHAR(255) UNIQUE NOT NULL,
  cost INTEGER NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE items;

-- +goose StatementEnd
