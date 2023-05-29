package request

type CreateUserRequest struct {
	OfficeID   int64  `json:"office_id" binding:"required,min=1"`
	UserName   string `json:"user_name" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email" binding:"required,email"`
}

type GetUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
