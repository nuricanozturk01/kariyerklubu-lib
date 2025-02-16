package response

import "net/http"

type Response[DATA any] struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	Data       DATA   `json:"data"`
}

func GetSuccessResponse[T any](message string, data T) *Response[T] {
	return &Response[T]{
		Message:    message,
		Success:    true,
		StatusCode: http.StatusOK,
		Data:       data,
	}
}

func GetResponse[T any](message string, statusCode int, data T) *Response[T] {
	return &Response[T]{
		Message:    message,
		Success:    statusCode == http.StatusOK,
		StatusCode: statusCode,
		Data:       data,
	}
}
