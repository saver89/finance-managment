-- name: CreateTransaction :one
insert into transactions (
  office_id, type, from_account_id, to_account_id, amount, currency_id, created_by
) values (
  $1, $2, $3, $4, $5, $6, $7
)
returning *;

-- name: GetTransaction :one
select * from transactions where id = $1 and deleted_at is null limit 1;

-- name: DeleteTransaction :exec
update transactions set deleted_at = now() where id = $1;