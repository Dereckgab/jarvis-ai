package http

// ErrorResponse is the standard error response body returned by all handlers.
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
