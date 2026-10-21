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

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeInitialized {
		panic("no status code set")
	}

	return rw.statusCode
}
