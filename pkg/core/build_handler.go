package iz2go

import (
	"slices"

	"github.com/gin-gonic/gin"
)

func BuildHandler(handler interface{}) *HandlerInfo {
	if handler == nil {
		return nil
	}
	if _, ok := handler.(IExecutable); !ok {
		panic("handler must implement IExecutable")
	}

	HandleInit(handler)
	method := ParseMethod(handler)
	handlerFunc := ParseHandler(handler)

	return &HandlerInfo{
		Method:  method,
		Handler: handlerFunc,
	}
}

func ParseMethod(handler interface{}) string {
	if h, ok := handler.(IWithMethod); ok {
		return h.GetMethod()
	}
	return "GET"
}

func HandleInit(handler interface{}) {
	if h, ok := handler.(IWithInit); ok {
		h.Init()
	}
}

func ParseHandler(handler interface{}) gin.HandlerFunc {
	if h, ok := handler.(IExecutableWithDecorator); ok {
		decorators := h.Decorators()
		slices.Reverse(decorators)
		handlerFunc := h.Execute
		for _, decorator := range decorators {
			handlerFunc = decorator(handlerFunc)
		}
		return handlerFunc
	}
	return handler.(IExecutable).Execute
}
