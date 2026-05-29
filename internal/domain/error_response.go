package domain

// ErrorResponse is the standard DIGIT error envelope.
type ErrorResponse struct {
	ResponseInfo ResponseInfo `json:"ResponseInfo"`
	Errors       []Error      `json:"Errors"`
}

// Error is a single DIGIT error item.
type Error struct {
	Code        string `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
	Params      string `json:"params,omitempty"`
}

// NewErrorResponse builds a failed DIGIT error response.
func NewErrorResponse(requestInfo *RequestInfo, errors []Error) ErrorResponse {
	return ErrorResponse{
		ResponseInfo: NewResponseInfoFromRequest(requestInfo, false),
		Errors:       errors,
	}
}
