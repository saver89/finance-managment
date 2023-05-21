package response

import db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"

type CreateTransferResponse struct {
	Transaction db.Transaction `json:"transaction"`
	FromAccount db.Account     `json:"from_account"`
	ToAccount   db.Account     `json:"to_account"`
}
