package request

type CreateAccountRequest struct {
	OfficeID   int64  `json:"office_id" db:"office_id"`
	Name       string `json:"name" db:"name" binding:"required"`
	CurrencyID int64  `json:"currency_id" db:"currency_id" binding:"required"`
	CreatedBy  int64  `json:"created_by" db:"created_by"`
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
