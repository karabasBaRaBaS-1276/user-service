package response

import "time"

type ErrorResponse struct {
	Timestamp string       `json:"timestamp"`
	Status    int          `json:"status"`
	Error     string       `json:"error"`
	Details   string       `json:"details"`
	Causes    []ErrorCause `json:"causes,omitempty"`
}

type ErrorCause struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

func NewError(
	status int,
	err string,
	details string,
	causes []ErrorCause,
) ErrorResponse {
	return ErrorResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Status:    status,
		Error:     err,
		Details:   details,
		Causes:    causes,
	}
}
