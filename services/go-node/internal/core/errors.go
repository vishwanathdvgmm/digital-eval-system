package core

import "fmt"

// Predefined error codes and helpers used across services.

const (
	CodeBadRequest      = "bad_request"
	CodeNotFound        = "not_found"
	CodeInternal        = "internal_error"
	CodeExternalService = "external_service_error"
	CodeUnauthorized    = "unauthorized"
)

// NewErrBody constructs a structured ErrBody.
func NewErrBody(code, message, detail string) *ErrBody {
	return &ErrBody{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// WrapErrorDetail formats an error with additional context for detail.
func WrapErrorDetail(err error, ctx string) string {
	if err == nil {
		return ctx
	}
	return fmt.Sprintf("%s: %v", ctx, err)
}
