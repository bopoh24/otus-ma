package helper

import (
	"errors"
	"net/http"
)

// ErrorResponse is a helper function to write error response
func ErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

var ErrInvalidClaims = errors.New("invalid claims")

type Claims struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
}

// ExtractClaims extracts claims from request headers
func ExtractClaims(r *http.Request) (Claims, error) {
	c := Claims{
		Id:        r.Header.Get("X-User"),
		Email:     r.Header.Get("X-Email"),
		FirstName: r.Header.Get("X-Given-Name"),
		LastName:  r.Header.Get("X-Family-Name"),
	}
	if c.Id == "" || c.Email == "" {
		return Claims{}, ErrInvalidClaims
	}
	return c, nil
}
