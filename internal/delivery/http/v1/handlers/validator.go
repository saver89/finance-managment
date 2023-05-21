package handlers

import "github.com/go-playground/validator/v10"

var ValidTransactionType validator.Func = func(fl validator.FieldLevel) bool {
	if transactionType, ok := fl.Field().Interface().(string); ok {
		return isValidTransactionType(transactionType)
	}

	return false
}

func isValidTransactionType(transactionType string) bool {
	switch transactionType {
	case "income", "outcome", "transfer", "adjustment":
		return true
	}

	return false
}
