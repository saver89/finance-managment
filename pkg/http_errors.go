package httpErrors

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/saver89/finance-management/pkg/logger"
)

const (
	BadRequest       = "Bad request"
	AlreadyExists    = "Already exists"
	NoSuchUser       = "User not found"
	WrongCredentials = "Wrong Credentials"
	NotFound         = "Not Found"
	Unauthorized     = "Unauthorized"
	Forbidden        = "Forbidden"
	BadQueryParams   = "Invalid query params"
	RequestTimeout   = "Request Timeout"
	InvalidEmail     = "Invalid email"
	InvalidPassword  = "Invalid password"
	InvalidField     = "Invalid field"
	ValidationError  = "Validation error"
)

var (
	ErrBadRequest            = errors.New("bad request")
	ErrWrongCredentials      = errors.New("wrong Credentials")
	ErrNotFound              = errors.New("not Found")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrPermissionDenied      = errors.New("permission Denied")
	ErrExpiredCSRFError      = errors.New("expired CSRF token")
	ErrWrongCSRFToken        = errors.New("wrong CSRF token")
	ErrCSRFNotPresented      = errors.New("cSRF not presented")
	ErrNotRequiredFields     = errors.New("no such required fields")
	ErrBadQueryParams        = errors.New("invalid query params")
	ErrInternalServerError   = errors.New("internal Server Error")
	ErrRequestTimeoutError   = errors.New("request Timeout")
	ErrExistsEmailError      = errors.New("user with given email already exists")
	ErrInvalidJWTToken       = errors.New("invalid JWT token")
	ErrInvalidJWTClaims      = errors.New("invalid JWT claims")
	ErrNotAllowedImageHeader = errors.New("not allowed image header")
	ErrNoCookie              = errors.New("not found cookie header")
	ErrInvalidUUID           = errors.New("invalid uuid")
	ErrValidationError       = errors.New("validation error")
)

// Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
	ErrBody() RestError
}

// Rest error struct
type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"err_causes,omitempty"`
}

// Error body
func (e RestError) ErrBody() RestError {
	return e
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

// Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// New Rest Error
func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New Rest Error From Bytes
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

// New Bad Request Error
func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  ErrBadRequest.Error(),
		ErrCauses: causes,
	}
}

// New Not Found Error
func NewNotFoundError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  ErrNotFound.Error(),
		ErrCauses: causes,
	}
}

// New Unauthorized Error
func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  ErrUnauthorized.Error(),
		ErrCauses: causes,
	}
}

// New Forbidden Error
func NewForbiddenError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  ErrForbidden.Error(),
		ErrCauses: causes,
	}
}

// New Internal Server Error
func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  ErrInternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

func NewValidationError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  ErrValidationError.Error(),
		ErrCauses: causes,
	}
	return result
}

// Parser of error string messages returns RestError
func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, NotFound, nil)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, RequestTimeout, nil)
	case errors.Is(err, ErrUnauthorized):
		return NewRestError(http.StatusUnauthorized, Unauthorized, nil)
	case errors.Is(err, ErrWrongCredentials):
		return NewRestError(http.StatusUnauthorized, Unauthorized, nil)
	case errors.Is(err, ErrValidationError):
		return NewRestError(http.StatusBadRequest, ValidationError, err.Error())
	default:
		if restErr, ok := err.(RestErr); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

// Error response
func ErrorResponse(log logger.Logger, err error) (int, interface{}) {
	log.Error(err.Error())
	return ParseErrors(err).Status(), ParseErrors(err)
}
