-- name: CreateAccount :one
insert into account (
  office_id, name, currency_id, created_by, state
) values (
  $1, $2, $3, $4, 'active'
) returning *;

-- name: GetAccount :one
select * from account where id = $1 and deleted_at is null;

-- name: GetAccountForUpdate :one
select * from account where id = $1 and deleted_at is null for no key update;

-- name: ListAccount :many
select * from account 
where deleted_at is null and office_id = $3
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
update account set name = $2, currency_id = $3
where id = $1 and deleted_at is null
returning *;

-- name: UpdateAccountBalance :one
update account set balance = $2
where id = $1 and deleted_at is null
returning *;

-- name: AddAccountBalance :one
update account set balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id) and deleted_at is null
returning *;

-- name: DeleteAccount :exec
update account set deleted_at = now() where id = $1;