-- name: CreateOfficeCurrencyRate :one
insert into office_currency_rate (
  office_id, from_currency_id, to_currency_id, rate
) values (
  $1, $2, $3, $4
)
returning *;

-- name: GetOfficeCurrencyRate :one
select * from office_currency_rate where id = $1 and deleted_at is null
LIMIT 1;

-- name: ListOfficeCurrencyRate :many
select * from office_currency_rate 
where deleted_at is null 
  and office_id = $1 
  and from_currency_id = $2 
  and to_currency_id = $3
LIMIT $4
OFFSET $5;

-- name: UpdateOfficeCurrencyRate :one
update office_currency_rate set rate = $2 where id = $1
returning *;

-- name: DeleteOfficeCurrencyRate :exec
update office_currency_rate set deleted_at = now() where id = $1;

