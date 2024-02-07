package base

type ApiResponse struct {
	Errors []ApiError `json:"errors"`
}

func NewApiResponse() ApiResponse {
	return ApiResponse{
		Errors: make([]ApiError, 0),
	}
}

func NewApiResponseWithError(errorMsgs ...string) ApiResponse {
	errors := make([]ApiError, len(errorMsgs))
	for i, errorMsg := range errorMsgs {
		errors[i] = ApiError{ErrorMsg: errorMsg}
	}
	return ApiResponse{
		Errors: errors,
	}
}

type ApiError struct {
	ErrorMsg string `json:"errorMessage"`
}
