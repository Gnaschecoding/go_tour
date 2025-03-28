package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	code   int      `json:"code"`
	msg    string   `json:"msg"`
	detail []string `json:"detail"`
}

// 管理 error的 code和信息 ，避免重复注册
var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code %d already exists", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d，错误信息： %s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}
func (e *Error) Details() []string {
	return e.detail
}
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.detail = []string{}

	for _, detail := range details {
		newError.detail = append(newError.detail, detail)
	}
	return &newError
}

func (e *Error) StatusCode() int {
	switch e.code {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
