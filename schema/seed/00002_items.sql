-- +goose Up
INSERT INTO
  items (codename, cost)
VALUES
  ('t-shirt', 80),
  ('cup', 20),
  ('book', 50),
  ('pen', 10),
  ('powerbank', 200),
  ('hoody', 300),
  ('umbrella', 200),
  ('socks', 10),
  ('wallet', 50) ON CONFLICT DO NOTHING;

;

-- +goose Down
DELETE FROM ITEMS;
