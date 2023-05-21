package request

type CreateTransferRequest struct {
	OfficeID      int64   `json:"office_id" db:"office_id" binding:"required,min=1"`
	FromAccountID int64   `json:"from_account_id" db:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" db:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" db:"amount" binding:"required,gt=0"`
	CurrencyID    int64   `json:"currency_id" db:"currency_id" binding:"required,min=1"`
	CreatedBy     int64   `json:"created_by" db:"created_by"`
	Type          string  `json:"type" db:"type" binding:"transaction-type"` // custom validator example
}
