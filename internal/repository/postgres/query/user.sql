-- name: CreateUser :one
insert into "user" (
  office_id, username, password_hash, first_name, last_name, 
  middle_name, birthdate, email, phone, created_by, state
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active')
returning *;

-- name: GetUser :one
select * from "user" where "id" = $1 and "deleted_at" is null;

-- name: ListUser :many
select * from "user" where "deleted_at" is null
LIMIT $1
OFFSET $2;;

-- name: UpdateUser :exec
update "user" set password_hash = $1, 
  first_name = $2, last_name = $3, middle_name = $4, birthdate = $5, 
  email = $6, phone = $7
where "id" = $8 and "deleted_at" is null;

-- name: DeleteUser :exec
update "user" set deleted_at = now() where "id" = $1;