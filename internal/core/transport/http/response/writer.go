package core_http_response

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	StatusCodeInitialized = -1
)

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeInitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeInitialized {
		return http.StatusOK
	}

	return rw.statusCode
}
