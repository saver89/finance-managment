// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: office_currency.sql

package db

import (
	"context"
)

const createOfficeCurrency = `-- name: CreateOfficeCurrency :one
INSERT INTO office_currency (
  office_id, currency_id, type
) VALUES (
  $1, $2, $3
)
returning id, office_id, currency_id, type, created_at, deleted_at
`

type CreateOfficeCurrencyParams struct {
	OfficeID   int64              `db:"office_id"`
	CurrencyID int64              `db:"currency_id"`
	Type       OfficeCurrencyType `db:"type"`
}

func (q *Queries) CreateOfficeCurrency(ctx context.Context, arg CreateOfficeCurrencyParams) (OfficeCurrency, error) {
	row := q.db.QueryRowContext(ctx, createOfficeCurrency, arg.OfficeID, arg.CurrencyID, arg.Type)
	var i OfficeCurrency
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.CurrencyID,
		&i.Type,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteOfficeCurrency = `-- name: DeleteOfficeCurrency :exec
UPDATE office_currency
  set deleted_at = now()
WHERE id = $1
`

func (q *Queries) DeleteOfficeCurrency(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteOfficeCurrency, id)
	return err
}

const getOfficeCurrency = `-- name: GetOfficeCurrency :one
SELECT id, office_id, currency_id, type, created_at, deleted_at FROM office_currency
WHERE id = $1 and deleted_at is null
LIMIT 1
`

func (q *Queries) GetOfficeCurrency(ctx context.Context, id int64) (OfficeCurrency, error) {
	row := q.db.QueryRowContext(ctx, getOfficeCurrency, id)
	var i OfficeCurrency
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.CurrencyID,
		&i.Type,
		&i.CreatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listOfficeCurrency = `-- name: ListOfficeCurrency :many
SELECT id, office_id, currency_id, type, created_at, deleted_at FROM office_currency
WHERE deleted_at is null and office_id = $1
`

func (q *Queries) ListOfficeCurrency(ctx context.Context, officeID int64) ([]OfficeCurrency, error) {
	rows, err := q.db.QueryContext(ctx, listOfficeCurrency, officeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OfficeCurrency
	for rows.Next() {
		var i OfficeCurrency
		if err := rows.Scan(
			&i.ID,
			&i.OfficeID,
			&i.CurrencyID,
			&i.Type,
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
