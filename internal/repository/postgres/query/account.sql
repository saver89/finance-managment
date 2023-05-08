-- name: CreateAccount :one
insert into account (
  office_id, name, currency_id, created_by, state
) values (
  $1, $2, $3, $4, 'active'
) returning *;

-- name: GetAccount :one
select * from account where id = $1 and deleted_at is null;

-- name: ListAccount :many
select * from account where deleted_at is null
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
update account set name = $2, currency_id = $3
where id = $1 and deleted_at is null
returning *;

-- name: DeleteAccount :exec
update account set deleted_at = now() where id = $1;