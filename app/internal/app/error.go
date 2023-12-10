package app

import "strconv"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) JSON() []byte {
	return []byte(`{"code":` + strconv.Itoa(e.Code) + `,"message":"` + e.Message + `"}`)
}
