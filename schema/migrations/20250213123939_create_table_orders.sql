-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id) ON DELETE RESTRICT,
  item_id INTEGER NOT NULL REFERENCES items (id) ON DELETE RESTRICT,
  transaction_id INTEGER NOT NULL REFERENCES transactions (id) ON DELETE RESTRICT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders (user_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;

-- +goose StatementEnd
