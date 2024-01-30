package helper

import (
	"net/http"
)

// ErrorResponse is a helper function to write error response
func ErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
