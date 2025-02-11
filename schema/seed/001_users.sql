-- +goose Up
INSERT INTO
  users (username, hashed_password, role)
VALUES
  ('employee', 'password', 'user'),
  ('director', 'password', 'admin');

-- +goose Down
DELETE FROM users
WHERE
  username in ('employee', 'director');
