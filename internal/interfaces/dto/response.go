package dto

type (
	Response struct {
		Code    int64       `json:"code"`
		Data    interface{} `json:"data"`
		Message string      `json:"msg"`
		Status  string      `json:"status"`
	}
	ResponseError struct {
		Code      int    `json:"code"`
		ErrorCode string `json:"error_code"`
		Message   string `json:"error_msg"`
	}
	NullResponse struct {
		Code    int64       `json:"code"`
		Data    interface{} `json:"data"`
		Message string      `json:"msg"`
	}
	ErrorCode struct {
		Data      interface{} `json:"data"`
		Code      int         `json:"http_code"`
		ErrorCode string      `json:"error_code"`
	}
)
