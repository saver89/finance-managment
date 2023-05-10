-- name: CreateHQOffice :one
INSERT INTO office (
  type, name, parent_id
) VALUES (
  'hq', $1, lastval()
)
RETURNING *;

-- name: AddOffice :one
INSERT INTO office (
  parent_id, type, name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOffice :one
SELECT * FROM office
WHERE id = $1 and deleted_at is null
LIMIT 1;

-- name: ListOffice :many
SELECT * FROM office
WHERE deleted_at is null and type = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;

-- name: UpdateOffice :one
UPDATE office
  set name = $2
WHERE id = $1 and deleted_at is null
RETURNING *;