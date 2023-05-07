-- name: CreateCurrency :one
INSERT INTO currency (
  name, short_name, state
) VALUES (
  $1, $2, 'active'
)
RETURNING *;

-- name: GetCurrency :one
SELECT * FROM currency
WHERE id = $1 LIMIT 1;

-- name: ListCurrency :many
SELECT * FROM currency
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateCurrency :one
UPDATE currency
  set name = $2,
  short_name = $3
WHERE id = $1
RETURNING *;
