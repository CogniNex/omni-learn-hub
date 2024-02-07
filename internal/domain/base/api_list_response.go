package base

type ApiListResponse[T any] struct {
	ApiResponse
	List []T
}

func NewApiListResponse[T any](list []T) ApiListResponse[T] {
	return ApiListResponse[T]{
		List: list,
	}
}

func NewApiListResponseWithError[T any](errorMsg string) ApiListResponse[T] {
	return ApiListResponse[T]{
		ApiResponse: NewApiResponseWithError(errorMsg),
	}
}
