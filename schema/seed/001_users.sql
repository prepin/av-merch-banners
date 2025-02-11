-- +goose Up
INSERT INTO
  users (username, hashed_password, role)
VALUES
  (
    'employee',
    '$2a$10$IMGLbeNyoaBT4xFC9qhN/.D3mks1wQ510baxLTI0Ie6zoMQR5ACEa',
    'user'
  ),
  (
    'director',
    '$2a$10$IMGLbeNyoaBT4xFC9qhN/.D3mks1wQ510baxLTI0Ie6zoMQR5ACEa',
    'admin'
  );

-- +goose Down
DELETE FROM users
WHERE
  username in ('employee', 'director');
