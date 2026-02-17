package response

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewError(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}
