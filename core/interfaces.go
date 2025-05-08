package core

import "github.com/gin-gonic/gin"

type IAPI interface {
	GetMethod() string
	Execute(c *gin.Context)
}

type IWithDecorator interface {
	UseDecorator(handler gin.HandlerFunc) gin.HandlerFunc
}
