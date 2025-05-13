package iz2go

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type HandlerInfo struct {
	Method   string
	Handler  gin.HandlerFunc
	ApiName  string
	Request  reflect.Type
	Response reflect.Type
}
