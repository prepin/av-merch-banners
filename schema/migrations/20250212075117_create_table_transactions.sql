-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users (id) ON DELETE RESTRICT,
  recipient_id INTEGER REFERENCES users (id) ON DELETE RESTRICT,
  amount INTEGER NOT NULL,
  transaction_type VARCHAR(50),
  transaction_reference_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_recipient_id ON transactions (recipient_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;

-- +goose StatementEnd
