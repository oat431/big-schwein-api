package common

type ResponseDTOError struct {
	HttpCode  int    `json:"http_code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
