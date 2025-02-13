-- +goose Up
-- +goose StatementBegin
INSERT INTO
  transactions (
    user_id,
    amount,
    transaction_type,
    transaction_reference_id
  )
SELECT
  id as user_id,
  1000 as amount,
  'CREDIT' as transaction_type,
  gen_random_uuid () as transaction_reference_id
FROM
  users;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM TRANSACTIONS;

-- +goose StatementEnd
