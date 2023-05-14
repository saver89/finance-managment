package response

type CreateAccountResponse struct {
	ID         int64  `json:"id"`
	OfficeID   int64  `db:"office_id"`
	Name       string `db:"name"`
	CurrencyID int64  `db:"currency_id"`
}
