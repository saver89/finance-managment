-- name: CreateOffice :one
INSERT INTO office (
  id, parent_id, type, name
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetOffice :one
SELECT * FROM office
WHERE id = $1 and deleted_at is null
LIMIT 1;

-- name: ListOffice :many
SELECT * FROM office
WHERE deleted_at is null
ORDER BY id DESC
LIMIT $1
OFFSET $2;

-- name: UpdateOffice :one
UPDATE office
  set name = $2,
  type = $3
WHERE id = $1 and deleted_at is null
RETURNING *;