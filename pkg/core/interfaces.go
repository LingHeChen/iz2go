package iz2go

import "github.com/gin-gonic/gin"

type Decorator = func(handler gin.HandlerFunc) gin.HandlerFunc

type IExecutable interface {
	Execute(c *gin.Context)
}

type IWithInit interface {
	Init()
}

type IWithMethod interface {
	GetMethod() string
}

type IWithDecorator interface {
	Decorators() []Decorator
}

type IExecutableWithDecorator interface {
	IExecutable
	IWithDecorator
}
