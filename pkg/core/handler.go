package iz2go

import "github.com/gin-gonic/gin"

type HandlerInfo struct {
	Method  string
	Handler gin.HandlerFunc
}
