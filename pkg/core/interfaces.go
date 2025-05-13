package iz2go

import "github.com/gin-gonic/gin"

type Decorator = func(handler gin.HandlerFunc) gin.HandlerFunc

type IWithInit interface {
	Init()
}

type IWithMethod interface {
	GetMethod() string
}

type IWithSummary interface {
	GetSummary() string
}

type IError interface {
	error
	GetCode() int
	GetMessage() string
}

type Error struct {
	Code    int
	Message string
}

func NewError(code int, message string) IError {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) GetMessage() string {
	return e.Message
}
