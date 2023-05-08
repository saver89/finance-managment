-- name: CreateOfficeCurrency :one
INSERT INTO office_currency (
  office_id, currency_id, type
) VALUES (
  $1, $2, $3
)
returning *;

-- name: GetOfficeCurrency :one
SELECT * FROM office_currency
WHERE id = $1 and deleted_at is null
LIMIT 1;

-- name: ListOfficeCurrency :many
SELECT * FROM office_currency
WHERE deleted_at is null and office_id = $1;

-- name: DeleteOfficeCurrency :exec
UPDATE office_currency
  set deleted_at = now()
WHERE id = $1;