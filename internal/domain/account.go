package domain

type Account struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Balance    float64 `json:"balance"`
	OfficeID   int64   `json:"office_id"`
	CurrencyID int64   `json:"currency_id"`
	CreatedBy  int64   `json:"created_by"`
	State      string  `json:"state"`
	CreatedAt  string  `json:"created_at"`
}
