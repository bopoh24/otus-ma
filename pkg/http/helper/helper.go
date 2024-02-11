package helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// ErrorResponse is a helper function to write error response
func ErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

// JSONResponse is a helper function to write JSON response
func JSONResponse(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
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

// JSONTime struct to make null for empty time
type JSONTime struct {
	time.Time
}

// MarshalJSON override method
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}
