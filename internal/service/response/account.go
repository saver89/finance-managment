package response

import "github.com/saver89/finance-management/internal/domain"

type CreateAccountResponse struct {
	ID         int64  `json:"id"`
	OfficeID   int64  `db:"office_id"`
	Name       string `db:"name"`
	CurrencyID int64  `db:"currency_id"`
}

type GetAccountResponse struct {
	Account domain.Account `json:"account"`
}

type ListAccountResponse struct {
	Accounts []domain.Account `json:"accounts"`
	Total    int64            `json:"total"`
}
