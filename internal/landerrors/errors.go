package landerrors

import "strings"

// CustomException mirrors org.egov.tracer.model.CustomException.
type CustomException struct {
	Code    string
	Message string
	Errors  map[string]string
}

func (e *CustomException) Error() string {
	if len(e.Errors) > 0 {
		parts := make([]string, 0, len(e.Errors))
		for k, v := range e.Errors {
			parts = append(parts, k+": "+v)
		}
		return strings.Join(parts, "; ")
	}
	if e.Message != "" {
		return e.Code + ": " + e.Message
	}
	return e.Code
}

func New(code, message string) *CustomException {
	return &CustomException{Code: code, Message: message}
}

func NewMap(errors map[string]string) *CustomException {
	return &CustomException{Errors: errors}
}
