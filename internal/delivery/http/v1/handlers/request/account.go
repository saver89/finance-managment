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

type ListAccountRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
	OfficeID int64 `form:"office_id" binding:"required,min=1"`
}
