package base

type ApiValueResponse struct {
	ApiResponse
	Value   interface{} `json:"value"`
	Success bool        `json:"success"`
}

func NewApiValueResponse(value interface{}) ApiValueResponse {
	return ApiValueResponse{
		Value:   value,
		Success: true,
	}
}

func NewApiValueResponseWithError(errorMsg string) ApiValueResponse {
	return ApiValueResponse{
		ApiResponse: NewApiResponseWithError(errorMsg),
		Success:     false,
	}
}
