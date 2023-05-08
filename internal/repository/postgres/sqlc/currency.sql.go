// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: currency.sql

package db

import (
	"context"
)

const createCurrency = `-- name: CreateCurrency :one
INSERT INTO currency (
  name, short_name, state
) VALUES (
  $1, $2, 'active'
)
RETURNING id, name, short_name, state, created_at, deleted_at
`

type CreateCurrencyParams struct {
	Name      string `db:"name"`
	ShortName string `db:"short_name"`
}

func (q *Queries) CreateCurrency(ctx context.Context, arg CreateCurrencyParams) (Currency, error) {
	row := q.db.QueryRowContext(ctx, createCurrency, arg.Name, arg.ShortName)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortName,
		&i.State,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteCurrency = `-- name: DeleteCurrency :exec
UPDATE currency
  set deleted_at = now()
WHERE id = $1
`

func (q *Queries) DeleteCurrency(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCurrency, id)
	return err
}

const getCurrency = `-- name: GetCurrency :one
SELECT id, name, short_name, state, created_at, deleted_at FROM currency
WHERE id = $1 and deleted_at is null
LIMIT 1
`

func (q *Queries) GetCurrency(ctx context.Context, id int64) (Currency, error) {
	row := q.db.QueryRowContext(ctx, getCurrency, id)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortName,
		&i.State,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listCurrency = `-- name: ListCurrency :many
SELECT id, name, short_name, state, created_at, deleted_at FROM currency
where deleted_at is null
ORDER BY short_name
LIMIT $1
OFFSET $2
`

type ListCurrencyParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) ListCurrency(ctx context.Context, arg ListCurrencyParams) ([]Currency, error) {
	rows, err := q.db.QueryContext(ctx, listCurrency, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Currency
	for rows.Next() {
		var i Currency
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ShortName,
			&i.State,
			&i.CreatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCurrency = `-- name: UpdateCurrency :one
UPDATE currency
  set name = $2,
  short_name = $3
WHERE id = $1
  and deleted_at is null
RETURNING id, name, short_name, state, created_at, deleted_at
`

type UpdateCurrencyParams struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	ShortName string `db:"short_name"`
}

func (q *Queries) UpdateCurrency(ctx context.Context, arg UpdateCurrencyParams) (Currency, error) {
	row := q.db.QueryRowContext(ctx, updateCurrency, arg.ID, arg.Name, arg.ShortName)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ShortName,
		&i.State,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}
