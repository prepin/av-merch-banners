-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id) ON DELETE RESTRICT,
  counterparty_id INTEGER REFERENCES users (id) ON DELETE RESTRICT,
  amount INTEGER NOT NULL,
  transaction_type VARCHAR(50),
  transaction_reference_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_counterparty_id ON transactions (counterparty_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;

-- +goose StatementEnd
