package response

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

func Success[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{Data: data}
}
