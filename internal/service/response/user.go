package response

import (
	db "github.com/saver89/finance-management/internal/repository/postgres/sqlc"
)

type CreateUserResponse struct {
	User db.User `json:"user"`
}

type GetUserResponse struct {
	User db.User `json:"user"`
}
