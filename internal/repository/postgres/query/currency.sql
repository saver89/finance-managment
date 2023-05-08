-- name: CreateCurrency :one
INSERT INTO currency (
  name, short_name, state
) VALUES (
  $1, $2, 'active'
)
RETURNING *;

-- name: GetCurrency :one
SELECT * FROM currency
WHERE id = $1 and deleted_at is null
LIMIT 1;

-- name: ListCurrency :many
SELECT * FROM currency
where deleted_at is null
ORDER BY short_name
LIMIT $1
OFFSET $2;

-- name: UpdateCurrency :one
UPDATE currency
  set name = $2,
  short_name = $3
WHERE id = $1
  and deleted_at is null
RETURNING *;

-- name: DeleteCurrency :exec
UPDATE currency
  set deleted_at = now()
WHERE id = $1;