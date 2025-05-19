package dto

type ApiResponse[T any] struct {
	ResponseKey     string `json:"status"`
	ResponseMessage string `json:"message"`
	Data            T      `json:"data"`
}

type ApiResponseSuccess[T any] struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   T      `json:"data"`
}

type ApiResponseFail struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
type ApiResponseSuccessNoData struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
